# Orcgen

Orcgen is a Go package that enables an easy-to-use conversion of web pages and HTML content to various file formats like PDF, PNG, and JPEG.
The underlying implementation uses the [Rod library](https://github.com/go-rod/rod) for the page conversion.

## Functionalities and packages

Orcgen provides the following functionalities:

- Conversion of web pages and HTML content to a static file (PNG, PDF...).

  - This can be done simply using the Generate function, but if you need
  prior configuration, you can access all the webdriver functionalities.
  - You can also use the other functions at orcgen.go, as specified in the examples page.

### Package folder

- FileInfo:
    A struct to standardize the returns and file saves.
    There's a Output function that writes the content to a output file.

- Handlers:
    The implementations of the page file save functionality (PDF / Screenshots).

- Webdriver:
    Simple wrapper over rod library.

## Installation

To use Orcgen, you can install it via Go modules:

```sh
    go get github.com/luabagg/orcgen
```

Then you can import it in your Go code:

```go
    import "github.com/luabagg/orcgen"
```

## Usage Example

The package comes with examples that demonstrate the usage of the various functions and features provided by Orcgen.
It's the way-to-go if you're trying to use this package for the first time.

```go
    import "github.com/luabagg/orcgen"

    orcgen.Generate(
        "https://www.github.com",
        proto.PageCaptureScreenshot{
            Format: proto.PageCaptureScreenshotFormatWebp,
        },
        "github.webp",
    )
```

You can more in [examples_test.go](https://github.com/luabagg/orcgen/tree/main/examples_test.go) page.

## Contributors

This project is an open-source project, and contributions from other developers are welcome. If you encounter any issues or have suggestions for improvement, please submit them on the project's GitHub page.
