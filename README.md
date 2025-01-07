<div align="center">

# http

![Build Status](https://img.shields.io/github/actions/workflow/status/fmich7/http/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/fmich7/http-server)](https://goreportcard.com/report/github.com/fmich7/http-server)
![Test Coverage](https://img.shields.io/badge/test--coverage-90%25-blue)

**HTTP server made from scratch in Go 🚀✨**

[Description](#📖-description) • [Features](#✨-features) • [Quick Start](#🚀-quick-start) • [Documentation](#💡-documentation) • [Testing](#🛠️-testing)

</div>

## 📖 Description

This project is a lightweight and easy to use **HTTP server** written in Go, built to follow the [RFC 7230](https://tools.ietf.org/html/rfc7230) (HTTP/1.1) standard. It's designed to make building web APIs and applications simple and straightforward while maintaining great performance and flexibility.

## ✨ Features

- **Supports HTTP methods**: Handle `GET`, `POST`, `PUT`, `DELETE`, and more.
- **Routing**: `Dynamic` and `static` route handling.
- **Built-in Timeouts**: Automatically `close slow connections` and prevent resource locks.
- **Request Timeout Handling**: Gracefully stop long-running requests by using `contexts`.
- **Custom Middleware Support**: Easily extend server’s functionality by adding reusable logic to handlers.
- **Ease of Usage**: `Start quickly` with `minimal setup` and easily add new features.
- **Test Coverage**: `90%+` of the code is tested for reliability.
- **RFC Compliance**: Compatibility with [HTTP/1.1](https://tools.ietf.org/html/rfc7230) standards.

## 🚀 Quick Start

### Get the package

```bash
go get github.com/fmich7/http
```

### Example

```go
package main

import (
	"fmt"

	"github.com/fmich7/http"
)

// Logging middleware
func LoggingMiddleware(next http.HTTPHandler) http.HTTPHandler {
	return func(r *http.HTTPRequest, w http.ResponseWriter) {
		fmt.Printf("Received %s request for %s\n", r.Method, r.URL)
		// Call the next handler
		next(r, w)
	}
}

// Handler function for the /hello/{who} route
// {who} is a dynamic route parameter
func HelloHandler(r *http.HTTPRequest, w http.ResponseWriter) {
	w.SetStatus(200)
	w.SetHeader("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<h1>Hello, %s!</h1>", r.Params["who"])))
}

func main() {
	router := http.NewHTTPRouter()
	server := http.NewServer(":3000", router)

	// Register a handler with middleware
	router.HandlerFunc("GET", "/hello/{who}", LoggingMiddleware(HelloHandler))

	// Start the server with error handling
	if err := server.Start(); err != nil {
		fmt.Println("Error occurred while starting the server:", err)
	}
}
```

## 💡 Documentation

Explore the full documentation for this package on

> [pkg.go.dev/github.com/fmich7/http](https://pkg.go.dev/github.com/fmich7/http#section-documentation)

## 🛠️ Testing

```bash
# go test -v -timeout 5s -race ./...
make test
# Generate test coverage file
make coverage
```
