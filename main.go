package main

import (
	"fmt"

	http "github.com/fmich7/http/pkg"
)

func main() {

	router := http.NewHTTPRouter()
	s := http.NewServer(":3000", router)

	router.HandlerFunc("GET", "/a/{asd}", func(r http.HTTPRequest, params map[string]string) http.HTTPResponse {
		content := []byte(fmt.Sprintln(params["asd"]))
		fmt.Println(r.Body)
		return http.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "text/plain", "Content-Length": fmt.Sprintf("%d", len(content))},
			Body:       content,
		}
	})
	router.HandlerFunc("GET", "/echo", func(h http.HTTPRequest, m map[string]string) http.HTTPResponse {
		return http.HTTPResponse{
			StatusCode: 200,
			Body:       []byte("Hello"),
		}
	})
	s.Start()
}
