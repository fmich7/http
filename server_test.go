package http

import (
	"fmt"
	"io"
	"net"
	"testing"
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
	s, port := startTestServer(t, NewHTTPRouter())
	defer s.Stop()

	// Connect to the server and send a message
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	expected := "HTTP/1.1 404 Not Found\r\nContent-Length: 10\r\nContent-Type: text/plain\r\n\r\nNot Found\n"
	request := "GET /unknown HTTP/1.1\r\n\r\n"
	// Write the request to the client side
	_, err = conn.Write([]byte(request))
	if err != nil {
		t.Fatalf("Failed to write to client connection: %s", err)
	}

	response, err := io.ReadAll(conn)
	if err != nil {
		t.Fatalf("failed to read: %s", err)
	}

	if string(response) != expected {
		t.Errorf("expected %q, but got %q", expected, response)
	}
}

func TestHandleConnection(t *testing.T) {
	s, port := startTestServer(t, NewHTTPRouter())
	defer s.Stop()

	// Register the handler before sending any requests
	s.router.HandlerFunc("GET", "/echo", func(r *HTTPRequest, w ResponseWriter) {
		w.Write([]byte("Hello\n"))
	})

	// Function to handle a connection and return the response
	sendRequest := func(request string) string {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		// Write the request to the client side
		_, err = conn.Write([]byte(request))
		if err != nil {
			t.Fatalf("Failed to write to client connection: %s", err)
		}

		response, err := io.ReadAll(conn)
		if err != nil {
			t.Fatalf("failed to read: %s", err)
		}

		return string(response)
	}

	tests := []struct {
		name string
		data string
		want string
	}{
		{
			name: "Bad request body",
			data: "INVALID REQUEST\r\n\r\n",
			want: "HTTP/1.1 400 Bad Request\r\nContent-Length: 12\r\nContent-Type: text/plain\r\n\r\nBad Request\n",
		},
		{
			name: "Invalid path",
			data: "GET /unknown HTTP/1.1\r\n\r\n",
			want: "HTTP/1.1 404 Not Found\r\nContent-Length: 10\r\nContent-Type: text/plain\r\n\r\nNot Found\n",
		},
		{
			name: "Registered path",
			data: "GET /echo HTTP/1.1\r\n\r\n",
			want: "HTTP/1.1 200 OK\r\nContent-Length: 6\r\nContent-Type: text/plain\r\n\r\nHello\n",
		},
	}

	for _, tt := range tests {
		got := sendRequest(tt.data)
		if tt.want != got {
			t.Errorf("Expected response %q, but got %q", tt.want, got)
		}
	}
}

func TestGetPort(t *testing.T) {
	s, port := startTestServer(t, NewHTTPRouter())
	defer s.Stop()

	expected := s.listener.Addr().(*net.TCPAddr).Port
	if port != expected {
		t.Errorf("Port is %d, got %d", expected, port)
	}
}

func TestServerSetAndIsRunning(t *testing.T) {
	server := &Server{}

	// Test initial state
	if server.isRunning() {
		t.Fatalf("Expected server to not be running, but got running=true")
	}

	// Test setting state to true
	server.setRunning(true)
	if !server.isRunning() {
		t.Fatalf("Expected server to be running, but got running=false")
	}

	// Test setting state to false
	server.setRunning(false)
	if server.isRunning() {
		t.Fatalf("Expected server to not be running, but got running=true")
	}
}
