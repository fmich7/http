package http

import (
	"bytes"
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

	var buf bytes.Buffer
	tmp := make([]byte, 1024)

	for {
		n, err := clientConn.Read(tmp)
		if err == io.EOF {
			break
		}

		if err != nil {
			t.Fatalf("read error: %s", err)
		}

		buf.Write(tmp[:n])
	}

	expected := "HTTP/1.1 200 OK\r\nContent-Length: 13\r\nContent-Type: text/plain\r\n\r\nHello, world!"
	got := buf.String()
	if got != expected {
		t.Errorf("Write() output = %q, want %q", got, expected)
	}
}
