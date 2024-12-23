package http

import (
	"fmt"
	"net"
	"sort"
)

// Status codes
const (
	StatusOK          = 200
	StatusBadRequest  = 400
	StatusNotFound    = 404
	StatusServerError = 500
)

// StatusDescription returns a status description for the given status code
func StatusDescription(code int) string {
	switch code {
	case 200:
		return "OK"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return ""
	}
}

type ResponseWriter interface {
	Write([]byte) (int, error)   // Write sends response to conn
	WriteHeader(statusCode int)  // Sets status code of response
	SetHeader(key, value string) // Sets headers of repsonse
}

// DefaultResponseWriter implements ResponseWriter interface
type DefaultResponseWriter struct {
	conn        net.Conn
	statusCode  int
	headers     map[string]string
	wroteHeader bool
}

// Returns new response writer
func NewResponseWriter(conn net.Conn) *DefaultResponseWriter {
	return &DefaultResponseWriter{
		conn:    conn,
		headers: make(map[string]string),
	}
}

// WriteHeader writes status code of response
func (rw *DefaultResponseWriter) WriteHeader(statusCode int) {
	if rw.wroteHeader {
		return
	}
	rw.statusCode = statusCode
	rw.wroteHeader = true
}

// SetHeader sets headers
func (rw *DefaultResponseWriter) SetHeader(key, value string) {
	rw.headers[key] = value
}

// Write sends resonse to client
func (rw *DefaultResponseWriter) Write(body []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(200)
	}

	if _, exists := rw.headers["Content-Type"]; !exists {
		rw.SetHeader("Content-Type", "text/plain")
	}

	if _, exists := rw.headers["Content-Length"]; !exists {
		rw.SetHeader("Content-Length", fmt.Sprintf("%d", len(body)))
	}

	responseLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", rw.statusCode, StatusDescription(rw.statusCode))
	if _, err := rw.conn.Write([]byte(responseLine)); err != nil {
		return 0, fmt.Errorf("failed to write response line: %w", err)
	}

	// Get the keys and sort them
	keys := make([]string, 0, len(rw.headers))
	for key := range rw.headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Write headers in sorted order
	for _, key := range keys {
		headerLine := fmt.Sprintf("%s: %s\r\n", key, rw.headers[key])
		if _, err := rw.conn.Write([]byte(headerLine)); err != nil {
			return 0, fmt.Errorf("failed to write header: %w", err)
		}
	}

	if _, err := rw.conn.Write([]byte("\r\n")); err != nil {
		return 0, fmt.Errorf("failed to write blank line after headers: %w", err)
	}

	n, err := rw.conn.Write(body)
	if err != nil {
		return n, fmt.Errorf("failed to write body: %w", err)
	}

	return n, nil
}
