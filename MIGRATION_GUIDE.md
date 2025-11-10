# Migration Guide: Moving to the New Type-Safe API

This guide helps you migrate from the deprecated generic `Generate()` and `NewHandler()` functions to the new type-safe API that eliminates runtime type assertions.

## Why Migrate?

The new API provides:
- **Compile-time type safety** - Catch errors before runtime
- **No type assertions** - Eliminates runtime type checking
- **Clearer intent** - Separate functions for HTML vs URL inputs
- **Better IDE support** - Improved autocomplete and documentation
- **Easier testing** - More straightforward to mock and test

## Quick Reference

| Old API | New API |
|---------|---------|
| `Generate(url, config, output)` | `GenerateURL(url, builder, output)` |
| `Generate(html, config, output)` | `GenerateHTML(html, builder, output)` |
| `NewHandler(config)` | `PDF(config)` or `Screenshot(config)` |

## Migration Examples

### Example 1: URL to Screenshot

**Before:**
```go
err := orcgen.Generate(
    "https://example.com",
    orcgen.ScreenshotConfig{
        Format: "png",
    },
    "output.png",
)
```

**After:**
```go
err := orcgen.GenerateURL(
    "https://example.com",
    orcgen.Screenshot(orcgen.ScreenshotConfig{
        Format: "png",
    }),
    "output.png",
)
```

### Example 2: HTML to PDF

**Before:**
```go
htmlBytes := []byte("<html>...</html>")
err := orcgen.Generate(
    htmlBytes,
    orcgen.PDFConfig{
        Landscape:       true,
        PrintBackground: true,
    },
    "output.pdf",
)
```

**After:**
```go
htmlBytes := []byte("<html>...</html>")
err := orcgen.GenerateHTML(
    htmlBytes,
    orcgen.PDF(orcgen.PDFConfig{
        Landscape:       true,
        PrintBackground: true,
    }),
    "output.pdf",
)
```

### Example 3: Using NewHandler with SetFullPage

**Before:**
```go
handler := orcgen.NewHandler(orcgen.PDFConfig{
    PrintBackground: true,
})
handler.SetFullPage(true)

fileinfo, err := orcgen.ConvertWebpage(handler, "https://example.com")
```

**After:**
```go
// Option 1: Using the builder pattern
err := orcgen.GenerateURL(
    "https://example.com",
    orcgen.PDF(orcgen.PDFConfig{
        PrintBackground: true,
    }).FullPage(true),
    "output.pdf",
)

// Option 2: If you need the handler directly
handler := pdf.New().SetConfig(orcgen.PDFConfig{
    PrintBackground: true,
}).SetFullPage(true)

wd := webdriver.FromDefault()
defer wd.Close()
page := wd.UrlToPage("https://example.com")
wd.WaitLoad(page)
fileinfo, err := handler.GenerateFile(page)
```

### Example 4: Advanced Usage with Custom WebDriver

**Before:**
```go
handler := orcgen.NewHandler(orcgen.ScreenshotConfig{
    Format: "webp",
})

wd := webdriver.FromDefault()
defer wd.Close()
page := wd.UrlToPage("https://example.com")
// ... custom page manipulation ...
wd.WaitLoad(page)

fileinfo, err := handler.GenerateFile(page)
fileinfo.Output("output.webp")
```

**After:**
```go
builder := orcgen.Screenshot(orcgen.ScreenshotConfig{
    Format: "webp",
})

wd := webdriver.FromDefault()
defer wd.Close()
page := wd.UrlToPage("https://example.com")
// ... custom page manipulation ...
wd.WaitLoad(page)

// Get the handler from the builder (internal method)
handler := screenshot.New().SetConfig(orcgen.ScreenshotConfig{
    Format: "webp",
})
fileinfo, err := handler.GenerateFile(page)
fileinfo.Output("output.webp")
```

## Breaking Changes

There are **no breaking changes** in this release. The old API is deprecated but still fully functional. You can migrate at your own pace.

### Deprecation Warnings

When you build your code, you may see deprecation warnings:
```
warning: orcgen.Generate is deprecated: Use GenerateHTML or GenerateURL instead for better type safety
warning: orcgen.NewHandler is deprecated: Use PDF() or Screenshot() factory functions instead
```

These are informational and won't prevent compilation.

## Step-by-Step Migration Process

1. **Identify all uses of `Generate()`**
   - Search your codebase for `orcgen.Generate(`
   - Determine if each call uses a URL (string) or HTML ([]byte)

2. **Replace with appropriate function**
   - String input → `GenerateURL()`
   - `[]byte` input → `GenerateHTML()`

3. **Update handler construction**
   - `NewHandler(PDFConfig{...})` → `PDF(PDFConfig{...})`
   - `NewHandler(ScreenshotConfig{...})` → `Screenshot(ScreenshotConfig{...})`

4. **Use method chaining for options**
   - Old: `handler := NewHandler(config); handler.SetFullPage(true)`
   - New: `PDF(config).FullPage(true)`

5. **Test your changes**
   - Run your test suite
   - Verify all generated files are correct

## Common Patterns

### Pattern: Simple URL to PDF
```go
// New API
orcgen.GenerateURL(
    url,
    orcgen.PDF(orcgen.PDFConfig{Landscape: true}),
    "output.pdf",
)
```

### Pattern: HTML with Full Page Screenshot
```go
// New API
orcgen.GenerateHTML(
    htmlBytes,
    orcgen.Screenshot(orcgen.ScreenshotConfig{
        Format: "png",
    }).FullPage(true),
    "output.png",
)
```

### Pattern: Chaining Multiple Options
```go
// New API supports fluent method chaining
orcgen.GenerateURL(
    url,
    orcgen.PDF(orcgen.PDFConfig{
        Landscape:       true,
        PrintBackground: true,
        PageRanges:      "1-5",
    }).FullPage(false),
    "output.pdf",
)
```

## Benefits You'll See

After migrating, you'll notice:

1. **Clearer code** - Intent is immediately obvious
2. **Better errors** - Type mismatches caught at compile time
3. **Faster development** - Better IDE autocomplete
4. **Easier debugging** - No mysterious type assertion panics
5. **More maintainable** - Less "magic" in the code

## Need Help?

If you encounter any issues during migration:
1. Check the examples in `examples_test.go`
2. Review the `ARCHITECTURE_ANALYSIS.md` for design rationale
3. Open an issue on GitHub

## Timeline

- **v2.x** - Old API deprecated but functional
- **v3.0** (future) - Old API may be removed (with advance notice)

We recommend migrating to the new API as soon as convenient to future-proof your code.
