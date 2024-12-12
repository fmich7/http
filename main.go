package main

import (
	"fmt"

	"github.com/fmich7/http-server/server"
)

func main() {

	router := server.NewHTTPRouter()
	s := server.NewServer(":3000", router)
	router.AddEndpoint("GET", "/", func(r server.HTTPRequest) server.HTTPResponse {
		content := []byte("ALL GOOD G")
		return server.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "text/plain", "Content-Length": fmt.Sprintf("%d", len(content))},
			Body:       content,
		}
	})
	s.Start()
}
