package server

import (
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	router := &HTTPRouter{}
	listenAddr := "localhost:8080"

	server := NewServer(listenAddr, router)

	if server.listenAddr != listenAddr {
		t.Errorf("Expected listenAddr to be %s, got %s", listenAddr, server.listenAddr)
	}

	if server.router != router {
		t.Error("Expected router to be the one passed to NewServer")
	}

	if server.quitch == nil {
		t.Error("Expected quitch channel to be initialized, got nil")
	}

	if server.listener != nil {
		t.Error("Expected listener to be nil initially")
	}
}

func TestStart(t *testing.T) {
	s := NewServer(":0", nil)

	go func() {
		err := s.Start()
		if err != nil {
			t.Errorf("Server failed to start: %s", err)
		}
	}()

	// Allow the server some time to start
	time.Sleep(10 * time.Millisecond)

	if s.listener == nil {
		t.Fatal("Expected server listener to be initialized, but got nil")
	}
}

func TestHandleConnection(t *testing.T) {
	s := NewServer(":0", NewHTTPRouter())

	// Register the handler before sending any requests
	s.router.HandlerFunc("GET", "/echo", func(h HTTPRequest, m map[string]string) HTTPResponse {
		return HTTPResponse{
			StatusCode: 200,
			Body:       []byte("Hello\n"),
		}
	})

	// Function to handle a connection and return the response
	sendRequest := func(request string) string {
		client, server := net.Pipe()
		defer client.Close()
		defer server.Close()

		var wg sync.WaitGroup
		wg.Add(1)

		// Run handleConnection in a goroutine
		go func() {
			defer wg.Done()
			s.handleConnection(server)
		}()

		// Write the request to the client side
		_, err := client.Write([]byte(request))
		if err != nil {
			t.Fatalf("Failed to write to client connection: %s", err)
		}

		// Read the response from the server side
		buf := make([]byte, 1024)
		n, err := client.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatalf("Failed to read from client connection: %s", err)
		}

		wg.Wait()
		return string(buf[:n])
	}

	tests := []struct {
		name string
		data string
		want string
	}{
		{"Bad request body", "INVALID REQUEST\r\n\r\n", "HTTP/1.1 400 Bad Request\r\n\r\nBad Request\n"},
		{"Invalid path", "GET /unknown HTTP/1.1\r\n\r\n", "HTTP/1.1 404 Not Found\r\n\r\nNot Found\n"},
		{"Registered path", "GET /echo HTTP/1.1\r\n\r\n", "HTTP/1.1 200 OK\r\n\r\nHello\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sendRequest(tt.data)
			if tt.want != got {
				t.Errorf("Expected response %s, but got %s", tt.want, got)
			}
		})
	}
}
