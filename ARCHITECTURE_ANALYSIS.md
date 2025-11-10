# Architecture Analysis & Improvement Recommendations

## Current Issues

### 1. Type Assertion-Based Strategy Pattern (Primary Issue)

**Location:** `orcgen.go:45-57` (NewHandler function)

```go
func NewHandler[Config handlers.Config](config Config) handlers.FileHandler[Config] {
    var handler any

    if _, ok := any(config).(PDFConfig); ok {
        handler = pdf.New()
    } else if _, ok := any(config).(ScreenshotConfig); ok {
        handler = screenshot.New()
    } else {
        panic("invalid config type provided")
    }

    return any(handler).(handlers.FileHandler[Config]).SetConfig(config)
}
```

**Problems:**
- Runtime type checking defeats the purpose of compile-time type safety
- Requires multiple unsafe type conversions through `any()`
- Panics on invalid types instead of returning errors
- Not easily extensible - adding new handlers requires modifying this function
- Violates Open/Closed Principle
- Makes testing harder
- Generic constraints don't actually prevent invalid types at compile time

### 2. Input Type Dispatch Issues

**Location:** `orcgen.go:30-34`

```go
if _, ok := any(html).([]byte); ok {
    fileinfo, err = ConvertHTML(handler, any(html).([]byte))
} else {
    fileinfo, err = ConvertWebpage(handler, any(html).(string))
}
```

**Problems:**
- Same type assertion pattern
- Assumes only two types are possible
- Duplicates the type assertion

### 3. Tight Coupling to Rod Proto Types

**Location:** `pkg/handlers/config.go:9-11`

```go
type Config interface {
    proto.PageCaptureScreenshot | proto.PagePrintToPDF
}
```

**Problems:**
- Handlers are tightly coupled to Rod's internal proto types
- Hard to mock for testing
- Leaks implementation details to the public API
- Cannot easily add custom handlers without modifying the constraint

### 4. Error Handling

- Uses `panic()` instead of returning errors
- No way to gracefully handle invalid configurations

---

## Recommended Solutions

### **Option 1: Functional Options Pattern (RECOMMENDED)**

This is the most idiomatic Go approach and provides excellent flexibility.

#### Design

```go
// Handler builder approach
type HandlerBuilder struct {
    handlerType string
    config      any
    fullPage    bool
}

func PDF(config PDFConfig) *HandlerBuilder {
    return &HandlerBuilder{
        handlerType: "pdf",
        config:      config,
    }
}

func Screenshot(config ScreenshotConfig) *HandlerBuilder {
    return &HandlerBuilder{
        handlerType: "screenshot",
        config:      config,
    }
}

func (hb *HandlerBuilder) FullPage(enabled bool) *HandlerBuilder {
    hb.fullPage = enabled
    return hb
}

// Usage
err := orcgen.GenerateHTML(htmlBytes,
    orcgen.PDF(orcgen.PDFConfig{
        Landscape: true,
        PrintBackground: true,
    }).FullPage(true),
    "output.pdf",
)

err := orcgen.GenerateURL("https://example.com",
    orcgen.Screenshot(orcgen.ScreenshotConfig{
        Format: "png",
    }),
    "output.png",
)
```

**Benefits:**
- No type assertions
- Clear, explicit API
- Compile-time type safety
- Easy to extend
- Self-documenting
- Chainable configuration

**Implementation Changes:**
- Remove generic Config constraint
- Remove NewHandler() function
- Each handler factory returns a concrete builder
- Separate GenerateHTML and GenerateURL functions

---

### **Option 2: Explicit Factory Functions**

Remove generics entirely and use explicit factory functions.

#### Design

```go
// Separate functions for each handler type
func GeneratePDF(input Input, config PDFConfig, output string) error
func GenerateScreenshot(input Input, config ScreenshotConfig, output string) error

// Or with handler instances
func NewPDFHandler(config PDFConfig) *PDFHandler
func NewScreenshotHandler(config ScreenshotConfig) *ScreenshotHandler

// Usage
handler := orcgen.NewPDFHandler(orcgen.PDFConfig{
    Landscape: true,
}).SetFullPage(true)

err := orcgen.GenerateHTML(htmlBytes, handler, "output.pdf")
```

**Benefits:**
- No type assertions
- No generics complexity
- Clear, explicit types
- Compile-time safety
- Simple to understand
- Easy to test

**Trade-offs:**
- More functions in the API
- Less "clever" but more maintainable

---

### **Option 3: Interface-Based Strategy Pattern**

Use traditional interface-based strategy pattern without generics.

#### Design

```go
// Base configuration interface
type HandlerConfig interface {
    Type() string
}

type PDFConfig struct {
    proto.PagePrintToPDF
}

func (c PDFConfig) Type() string { return "pdf" }

type ScreenshotConfig struct {
    proto.PageCaptureScreenshot
}

func (c ScreenshotConfig) Type() string { return "screenshot" }

// Handler interface without generics
type FileHandler interface {
    GenerateFile(page *rod.Page) (*fileinfo.Fileinfo, error)
    SetFullPage(bool) FileHandler
}

// Factory registry
type HandlerFactory func() FileHandler

var handlerRegistry = map[string]HandlerFactory{
    "pdf":        func() FileHandler { return pdf.New() },
    "screenshot": func() FileHandler { return screenshot.New() },
}

func NewHandler(config HandlerConfig) (FileHandler, error) {
    factory, ok := handlerRegistry[config.Type()]
    if !ok {
        return nil, fmt.Errorf("unknown handler type: %s", config.Type())
    }

    handler := factory()
    // Apply config through a separate method
    return handler.ApplyConfig(config), nil
}

// Usage
config := orcgen.PDFConfig{
    PagePrintToPDF: proto.PagePrintToPDF{
        Landscape: true,
    },
}

handler, err := orcgen.NewHandler(config)
if err != nil {
    return err
}

err = orcgen.GenerateHTML(htmlBytes, handler, "output.pdf")
```

**Benefits:**
- No type assertions for handler selection
- Registry pattern allows external registration
- Returns errors instead of panicking
- Traditional Go patterns
- Easy to extend without modifying core code

**Trade-offs:**
- Still need to handle config type somehow
- Slightly more complex setup

---

### **Option 4: Type-Safe Builder with Method Chaining**

Create a fluent API that builds the correct handler.

#### Design

```go
// Start with a generator that hasn't chosen a handler
type Generator struct {
    input  Input
    output string
}

func FromHTML(html []byte) *Generator {
    return &Generator{input: htmlInput{html}}
}

func FromURL(url string) *Generator {
    return &Generator{input: urlInput{url}}
}

// Method to select PDF output
func (g *Generator) ToPDF(config PDFConfig) *PDFGenerator {
    return &PDFGenerator{
        generator: g,
        config:    config,
    }
}

// Method to select Screenshot output
func (g *Generator) ToScreenshot(config ScreenshotConfig) *ScreenshotGenerator {
    return &ScreenshotGenerator{
        generator: g,
        config:    config,
    }
}

type PDFGenerator struct {
    generator *Generator
    config    PDFConfig
    fullPage  bool
}

func (pg *PDFGenerator) FullPage(enabled bool) *PDFGenerator {
    pg.fullPage = enabled
    return pg
}

func (pg *PDFGenerator) Save(output string) error {
    // Implementation
}

// Usage
err := orcgen.FromHTML(htmlBytes).
    ToPDF(orcgen.PDFConfig{
        Landscape: true,
    }).
    FullPage(true).
    Save("output.pdf")

err := orcgen.FromURL("https://example.com").
    ToScreenshot(orcgen.ScreenshotConfig{
        Format: "png",
    }).
    Save("output.png")
```

**Benefits:**
- No type assertions
- Type-safe at every step
- Fluent, intuitive API
- Self-documenting
- Each method returns the correct type
- Compile-time guarantees

**Trade-offs:**
- More types to maintain
- Slightly more boilerplate

---

## Addressing the Input Type Issue

All solutions should also fix the `string | []byte` type assertion problem.

### Recommended Approach: Separate Functions

```go
func GenerateHTML(html []byte, handler HandlerBuilder, output string) error
func GenerateURL(url string, handler HandlerBuilder, output string) error
```

**Benefits:**
- Clear intent
- No runtime type checking
- Better API documentation
- Compile-time safety

---

## Additional Improvements

### 1. Return Errors Instead of Panicking

```go
// Current
func NewHandler[Config handlers.Config](config Config) handlers.FileHandler[Config] {
    // ... panics on invalid type
}

// Better
func NewHandler(config Config) (FileHandler, error) {
    // ... returns error on invalid type
}
```

### 2. Decouple from Rod Proto Types

Create your own configuration types that wrap Rod's types:

```go
type PDFOptions struct {
    Landscape       bool
    PrintBackground bool
    PageRanges      string
    // ... other options
}

func (opts PDFOptions) toProto() proto.PagePrintToPDF {
    return proto.PagePrintToPDF{
        Landscape:       opts.Landscape,
        PrintBackground: opts.PrintBackground,
        PageRanges:      opts.PageRanges,
    }
}
```

**Benefits:**
- Public API not tied to Rod
- Can change underlying implementation
- Easier to test
- More control over API surface

### 3. Use Interfaces for Testing

```go
type PageRenderer interface {
    PDF(*proto.PagePrintToPDF) (io.Reader, error)
    Screenshot(bool, *proto.PageCaptureScreenshot) ([]byte, error)
}

// rod.Page already implements this, but now you can mock it
```

### 4. Error Types

```go
type InvalidHandlerError struct {
    Type string
}

func (e *InvalidHandlerError) Error() string {
    return fmt.Sprintf("invalid handler type: %s", e.Type)
}
```

---

## Recommended Migration Path

**Phase 1: Immediate Improvements (Backward Compatible)**
1. Add separate `GenerateHTML` and `GenerateURL` functions
2. Keep the generic `Generate` function for backward compatibility but mark it as deprecated
3. Return errors instead of panics in NewHandler

**Phase 2: Introduce Better API (New Major Version)**
1. Implement **Option 1 (Functional Options)** or **Option 4 (Builder Pattern)**
2. Remove generic constraints
3. Create explicit handler factories
4. Separate input methods

**Phase 3: Decouple Dependencies**
1. Create own config types
2. Wrap Rod proto types
3. Add proper interfaces for testing

---

## Comparison Matrix

| Approach | Type Safety | Extensibility | API Clarity | Complexity | Testability |
|----------|-------------|---------------|-------------|------------|-------------|
| **Current (Type Assertions)** | Runtime only | Low | Medium | High | Low |
| **Option 1: Functional Options** | Compile-time | High | Excellent | Low | Excellent |
| **Option 2: Explicit Factories** | Compile-time | Medium | Excellent | Very Low | Excellent |
| **Option 3: Interface Registry** | Partial | Very High | Good | Medium | Good |
| **Option 4: Type-Safe Builder** | Compile-time | Medium | Excellent | Medium | Excellent |

---

## My Recommendation: Option 1 (Functional Options Pattern)

This is the most idiomatic Go approach and provides:
- **No type assertions** - completely eliminated
- **Compile-time safety** - catch errors before runtime
- **Clear API** - self-documenting and intuitive
- **Easy to extend** - add new handlers without modifying core code
- **Excellent testing** - easy to mock and test
- **Familiar pattern** - commonly used in Go libraries

### Example of Proposed API

```go
// PDF from HTML
err := orcgen.GenerateHTML(
    htmlBytes,
    orcgen.PDF(orcgen.PDFConfig{
        Landscape:       true,
        PrintBackground: true,
    }).FullPage(true),
    "output.pdf",
)

// Screenshot from URL
err := orcgen.GenerateURL(
    "https://example.com",
    orcgen.Screenshot(orcgen.ScreenshotConfig{
        Format:  "png",
        Quality: proto.Int(90),
    }),
    "screenshot.png",
)

// Advanced usage with webdriver customization
wd := orcgen.NewWebDriver(orcgen.WebDriverConfig{
    LoadTimeout:  30 * time.Second,
    PageIdleTime: 2 * time.Second,
})
defer wd.Close()

fileInfo, err := orcgen.ConvertWithDriver(
    wd,
    "https://example.com",
    orcgen.PDF(orcgen.PDFConfig{
        Landscape: true,
    }),
)
```

---

## Questions to Consider

1. **Backward Compatibility:** Do you need to maintain the current API, or is this a new major version?
2. **Extensibility:** Do you want users to be able to add custom handlers, or keep it internal only?
3. **Complexity:** How simple should the API be for the 80% use case?
4. **Testing:** How important is it to mock Rod dependencies?

Let me know which direction you'd like to take, and I can help implement it!
