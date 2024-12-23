package main

import (
	"fmt"
	"os"

	http "github.com/fmich7/http/pkg"
)

func main() {
	router := http.NewHTTPRouter()
	s := http.NewServer(":3000", router)

	router.HandlerFunc("GET", "/echo/{asd}", func(h http.HTTPRequest, w http.ResponseWriter, m map[string]string) {
		w.Write([]byte(m["asd"]))
	})

	// Download a file
	router.HandlerFunc("GET", "/static/{file}", func(r http.HTTPRequest, w http.ResponseWriter, m map[string]string) {
		file, err := os.ReadFile("static/" + m["file"])
		fmt.Println(m["file"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(http.StatusDescription(http.StatusNotFound)))
			return
		}

		w.SetHeader("Content-Type", "application/octet-stream")
		w.Write([]byte(file))
	})

	// Post file
	router.HandlerFunc("POST", "/files/{file}", func(r http.HTTPRequest, w http.ResponseWriter, m map[string]string) {
		filename := m["file"]
		file, err := os.Create("upload/" + filename)
		defer file.Close()
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(http.StatusDescription(http.StatusNotFound)))
			return
		}

		n, err := file.Write(r.Body)
		if err != nil || n < len(r.Body) {
			w.WriteHeader(http.StatusServerError)
			w.Write([]byte(http.StatusDescription(http.StatusServerError)))
			return
		}
		w.Write([]byte("Successfully uploaded file"))
	})

	s.Start()
}
