package main

import (
	"fmt"

	"github.com/fmich7/http-server/server"
)

func main() {

	router := server.NewHTTPRouter()
	s := server.NewServer(":3000", router)

	router.HandlerFunc("GET", "/a/{asd}", func(r server.HTTPRequest, params map[string]string) server.HTTPResponse {
		content := []byte(fmt.Sprintln(params["asd"]))
		fmt.Println(r.Body)
		return server.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "text/plain", "Content-Length": fmt.Sprintf("%d", len(content))},
			Body:       content,
		}
	})
	router.HandlerFunc("GET", "/echo", func(h server.HTTPRequest, m map[string]string) server.HTTPResponse {
		return server.HTTPResponse{
			StatusCode: 200,
			Body:       []byte("Hello"),
		}
	})
	s.Start()
}
