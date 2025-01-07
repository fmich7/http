package http

import (
	"io"
	"net"
	"testing"
	"time"
)

const writeTimeout = 5 * time.Second

func TestStatusDescription(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{200, "OK"},
		{201, "Created"},
		{400, "Bad Request"},
		{404, "Not Found"},
		{408, "Request Timeout"},
		{500, "Internal Server Error"},
		{504, "Gateway Timeout"},
		{1337, ""},
	}

	for _, v := range tests {
		if got := StatusDescription(v.input); got != v.expected {
			t.Errorf("Failed test [%d], want %s got %s\n", v.input, v.expected, got)
		}
	}
}

func TestResponseWriter_Write(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	go func() {
		defer serverConn.Close()
		rw := NewResponseWriter(serverConn, writeTimeout)
		if _, err := rw.Write([]byte("Hello, world!")); err != nil {
			t.Errorf("Write error: %s", err)
		}
	}()

	// Read response
	in, err := io.ReadAll(clientConn)
	if err != nil {
		t.Fatalf("failed to read: %s", err)
	}

	expected := "HTTP/1.1 200 OK\r\nContent-Length: 13\r\nContent-Type: text/plain\r\n\r\nHello, world!"
	got := string(in)
	if got != expected {
		t.Errorf("Write() output = %q, want %q", got, expected)
	}
}
