package main

import (
	"fmt"
	"os"
	"time"

	http "github.com/fmich7/http/pkg"
)

func AsdHandler(r *http.HTTPRequest, w http.ResponseWriter) {
	time.Sleep(1 * time.Second)
	w.Write([]byte(r.Params["asd"]))
}

func main() {
	router := http.NewHTTPRouter()
	s := http.NewServer(":3000", router)

	// Params with timeout
	router.HandlerFunc("GET", "/echo/{asd}", http.TimeoutHandler(AsdHandler, 5*time.Second))

	// Download a file
	router.HandlerFunc("GET", "/static/{file}", func(r *http.HTTPRequest, w http.ResponseWriter) {
		file, err := os.ReadFile("static/" + r.Params["file"])
		fmt.Println(r.Params["file"])
		if err != nil {
			w.SetStatus(http.StatusNotFound)
			w.Write([]byte(http.StatusDescription(http.StatusNotFound)))
			return
		}

		w.SetHeader("Content-Type", "application/octet-stream")
		w.Write([]byte(file))
	})

	// Post file
	router.HandlerFunc("POST", "/files/{file}", func(r *http.HTTPRequest, w http.ResponseWriter) {
		filename := r.Params["file"]
		file, err := os.Create("upload/" + filename)
		defer file.Close()
		if err != nil {
			w.SetStatus(http.StatusNotFound)
			w.Write([]byte(http.StatusDescription(http.StatusNotFound)))
			return
		}

		n, err := file.Write(r.Body)
		if err != nil || n < len(r.Body) {
			w.SetStatus(http.StatusServerError)
			w.Write([]byte(http.StatusDescription(http.StatusServerError)))
			return
		}
		w.Write([]byte("Successfully uploaded file"))
	})

	s.Start()
}
