package http

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"
)

func TestTimeoutHandler(t *testing.T) {
	// Initial setup
	router := NewHTTPRouter()
	timeoutDuration := 20 * time.Millisecond

	// Define a handler that intentionally takes too long
	tmpHandler := func(r *HTTPRequest, w ResponseWriter) {
		time.Sleep(50 * time.Millisecond)
		w.Write([]byte("OK"))
	}
	router.HandlerFunc("GET", "/timeout", TimeoutHandler(tmpHandler, timeoutDuration))

	// Define a handler that completes on time
	fastHandler := func(r *HTTPRequest, w ResponseWriter) {
		w.Write([]byte("Fast Response"))
	}
	router.HandlerFunc("GET", "/ok", TimeoutHandler(fastHandler, timeoutDuration))

	s := NewServer(":0", router)

	// Helper function to simulate client-server interaction
	sendRequest := func(request string) string {
		client, server := net.Pipe()
		defer client.Close()
		defer server.Close()

		// Handle the connection in a separate goroutine
		go func() {
			s.handleConnection(server)
		}()

		// Send the request
		_, err := client.Write([]byte(request))
		if err != nil {
			t.Fatalf("Failed to write to client connection: %s", err)
		}

		// Read the response
		var buf bytes.Buffer
		tmp := make([]byte, 1024)
		for {
			n, err := client.Read(tmp)
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("Failed to read from client connection: %s", err)
			}
			buf.Write(tmp[:n])
		}

		return buf.String()
	}

	// Test cases
	tests := []struct {
		name         string
		request      string
		expectedBody string
		expectedCode string
	}{
		{
			name:         "Request times out",
			request:      "GET /timeout HTTP/1.1\r\n\r\n",
			expectedBody: "Gateway Timeout",
			expectedCode: "504",
		},
		{
			name:         "Request succeeds",
			request:      "GET /ok HTTP/1.1\r\n\r\n",
			expectedBody: "Fast Response",
			expectedCode: "200",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sendRequest(tt.request)
			if !bytes.Contains([]byte(got), []byte(tt.expectedCode)) {
				t.Errorf("Expected status code %s, got response: %s", tt.expectedCode, got)
			}
			if !bytes.Contains([]byte(got), []byte(tt.expectedBody)) {
				t.Errorf("Expected body %q, got response: %s", tt.expectedBody, got)
			}
		})
	}
}
