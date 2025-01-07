package http

import (
	"bytes"
	"fmt"
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
		time.Sleep(5000 * time.Millisecond)
		w.Write([]byte("OK"))
	}

	// Define a handler that completes on time
	fastHandler := func(r *HTTPRequest, w ResponseWriter) {
		w.Write([]byte("Fast Response"))
	}

	router.HandlerFunc("GET", "/timeout", TimeoutHandler(tmpHandler, timeoutDuration))
	router.HandlerFunc("GET", "/ok", TimeoutHandler(fastHandler, timeoutDuration))
	// ---------------------------

	s, port := startTestServer(t, router)
	defer s.Stop()

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
		got := sendRequest(tt.request)
		if !bytes.Contains([]byte(got), []byte(tt.expectedCode)) {
			t.Errorf("Expected status code %s, got response: %s", tt.expectedCode, got)
		}
		if !bytes.Contains([]byte(got), []byte(tt.expectedBody)) {
			t.Errorf("Expected body %q, got response: %s", tt.expectedBody, got)
		}
	}
}
