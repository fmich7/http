package main

import (
	"fmt"
	"os"
	"strconv"

	http "github.com/fmich7/http/pkg"
)

func main() {
	router := http.NewHTTPRouter()
	s := http.NewServer(":3000", router)

	// Download a file
	router.HandlerFunc("GET", "/static/{file}", func(h http.HTTPRequest, m map[string]string) http.HTTPResponse {
		file, err := os.ReadFile("static/" + m["file"])
		fmt.Println(m["file"])
		if err != nil {
			return http.HTTPResponse{
				StatusCode: http.StatusNotFound,
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
				Body: []byte("Not Found"),
			}
		}
		return http.HTTPResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type":   "application/octet-stream",
				"Content-Length": strconv.Itoa(len(file)),
			},
			Body: file,
		}
	})

	// Post file
	router.HandlerFunc("POST", "/files/{file}", func(h http.HTTPRequest, m map[string]string) http.HTTPResponse {
		filename := m["file"]
		file, err := os.Create("upload/" + filename)
		defer file.Close()
		if err != nil {
			return http.HTTPResponse{
				StatusCode: http.StatusNotFound,
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
				Body: []byte("Not Found"),
			}
		}
		fmt.Println(h.Body)
		n, err := file.Write(h.Body)
		if err != nil || n < len(h.Body) {
			return http.HTTPResponse{
				StatusCode: http.StatusServerError,
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
				Body: []byte("Internal server error"),
			}
		}

		return http.HTTPResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
			Body: []byte("Successfully uploaded file"),
		}
	})

	s.Start()
}
