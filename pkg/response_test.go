package http

import (
	"net"
	"testing"
)

func TestStatusDescription(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{200, "OK"},
		{404, "Not Found"},
		{400, "Bad Request"},
		{500, "Internal Server Error"},
		{1337, ""},
	}

	for _, v := range tests {
		if got := StatusDescription(v.input); got != v.expected {
			t.Errorf("Failed test [%d], want %s got %s\n", v.input, v.expected, got)
		}
	}
}

func TestHTTPResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response HTTPResponse
		expected string
	}{
		{
			name: "200 OK with headers and body",
			response: HTTPResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Type": "text/html",
					"Connection":   "close",
				},
				Body: []byte("<html><body>Hello, world!</body></html>"),
			},
			expected: "HTTP/1.1 200 OK\r\n" +
				"Content-Type: text/html\r\n" +
				"Connection: close\r\n" +
				"\r\n" +
				"<html><body>Hello, world!</body></html>",
		},
		{
			name: "404 Not Found with no headers and empty body",
			response: HTTPResponse{
				StatusCode: 404,
				Headers:    map[string]string{},
				Body:       []byte{},
			},
			expected: "HTTP/1.1 404 Not Found\r\n" +
				"\r\n",
		},
		{
			name: "500 Internal Server Error with headers",
			response: HTTPResponse{
				StatusCode: 500,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: []byte(`{"error": "server failure"}`),
			},
			expected: "HTTP/1.1 500 Internal Server Error\r\n" +
				"Content-Type: application/json\r\n" +
				"\r\n" +
				`{"error": "server failure"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.response.String()
			if got != tt.expected {
				t.Errorf("HTTPResponse.String() =\n%s\nWant:\n%s", got, tt.expected)
			}
		})
	}
}

func TestHTTPResponse_Write(t *testing.T) {
	response := HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
		Body: []byte("Hello, world!"),
	}

	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	go func() {
		err := response.Write(serverConn)
		if err != nil {
			t.Errorf("Write() error = %v", err)
		}
	}()

	buf := make([]byte, 2048)
	n, err := clientConn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read from conn: %v", err)
	}
	buf = buf[:n]

	expected := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, world!"
	if string(buf) != expected {
		t.Errorf("Write() output = %q, want %q", string(buf), expected)
	}
}
