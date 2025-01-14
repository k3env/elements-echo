# Stoplight Elements Middleware for Echo

Middleware for [Labstack Echo](https://github.com/labstack/echo) that integrates [Stoplight Elements](https://github.com/stoplightio/elements) to provide an interactive OpenAPI documentation interface.

# ⚠️ Project Status

This project is in early development, has no stable release, and is not ready for production use.

## Features
- Serve OpenAPI documentation using Stoplight Elements
- Easy integration with Echo applications
- Supports both JSON and YAML OpenAPI specifications

## Installation

```bash
go get github.com/k3env/elements-echo
```

## Usage

With embedded spec

```go
package main

import (
	"embed"
	"github.com/labstack/echo/v4"
	elements "github.com/k3env/elements-echo"
)

//go:embed all:docs
var embedFs embed.FS

func main() {
	e := echo.New()
	
	middleware := elements.New("/openapi")
	middleware, err := middleware.UseEmbed(embedFs, "docs/swagger.yaml")
	
	if err != nil {
		e.Logger.Fatal(err)
    }
	
	// Register the middleware to serve OpenAPI documentation
	e.Use(middleware.Handle())
	
	e.Logger.Fatal(e.Start(":8080"))
}
```

With spec from fs

```go
package main

import (
	"github.com/labstack/echo/v4"
	elements "github.com/k3env/elements-echo"
)

func main() {
	e := echo.New()
	
	middleware := elements.New("/openapi")
	middleware, err := middleware.UseSpecFile("docs/swagger.yaml")
	
	if err != nil {
		e.Logger.Fatal(err)
    }
	
	// Register the middleware to serve OpenAPI documentation
	e.Use(middleware.Handle())
	
	e.Logger.Fatal(e.Start(":8080"))
}
```

## API

### `New(urlPrefix string) echo.MiddlewareFunc`
Initializes the middleware with supplied route.

## License

[Apache 2.0](LICENSE)

