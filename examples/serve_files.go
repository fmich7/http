package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fmich7/http"
)

func SomeHandler(r *http.HTTPRequest, w http.ResponseWriter) {
	time.Sleep(100 * time.Second)
	w.Write([]byte(r.Params["asd"]))
}

func main2() {
	router := http.NewHTTPRouter()
	s := http.NewServer(":3000", router)

	// Register a handler
	router.HandlerFunc("GET", "/hello/{who}", func(r *http.HTTPRequest, w http.ResponseWriter) {
		w.SetStatus(200)
		w.SetHeader("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf("<h1>Hello, %s!</h1>", r.Params["who"])))
	})

	// Params with timeout
	router.HandlerFunc("GET", "/echo/{asd}", http.TimeoutHandler(SomeHandler, 100*time.Second))

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
		w.SetStatus(201)
		w.Write([]byte("Successfully uploaded file"))
	})

	if err := s.Start(); err != nil {
		log.Println("Error while starting server:", err)
	}
}
