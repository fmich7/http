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
