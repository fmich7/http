package http

import (
	"fmt"
	"net"
	"sort"
	"time"
)

// Status codes
const (
	StatusOK             = 200
	StatusBadRequest     = 400
	StatusNotFound       = 404
	StatusRequestTimeout = 408
	StatusServerError    = 500
	StatusGatewayTimeout = 504
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
	case 408:
		return "Request Timeout"
	case 500:
		return "Internal Server Error"
	case 504:
		return "Gateway Timeout"
	default:
		return ""
	}
}

type ResponseWriter interface {
	Write([]byte) (int, error)   // Write sends response to conn
	SetStatus(statusCode int)    // Sets status code of response
	SetHeader(key, value string) // Sets headers of repsonse
}

// DefaultResponseWriter implements ResponseWriter interface
type DefaultResponseWriter struct {
	conn         net.Conn
	statusCode   int
	headers      map[string]string
	wroteStatus  bool
	writeTimeout time.Duration
}

// Returns new response writer
func NewResponseWriter(conn net.Conn, writeTimeout time.Duration) *DefaultResponseWriter {
	return &DefaultResponseWriter{
		conn:         conn,
		headers:      make(map[string]string),
		writeTimeout: writeTimeout,
	}
}

// SetStatus writes status code of response
func (rw *DefaultResponseWriter) SetStatus(statusCode int) {
	if rw.wroteStatus {
		return
	}
	rw.statusCode = statusCode
	rw.wroteStatus = true
}

// SetHeader sets headers
func (rw *DefaultResponseWriter) SetHeader(key, value string) {
	rw.headers[key] = value
}

// Write sends resonse to client
func (rw *DefaultResponseWriter) Write(body []byte) (int, error) {
	// Set default status
	if !rw.wroteStatus {
		rw.SetStatus(200)
	}

	// Set default content type
	if _, exists := rw.headers["Content-Type"]; !exists {
		rw.SetHeader("Content-Type", "text/plain")
	}

	// Set Content-Length
	if _, exists := rw.headers["Content-Length"]; !exists {
		rw.SetHeader("Content-Length", fmt.Sprintf("%d", len(body)))
	}

	// Set write deadlines
	rw.conn.SetWriteDeadline(time.Now().Add(rw.writeTimeout))
	defer rw.conn.SetWriteDeadline(time.Time{})

	// Write status
	responseLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", rw.statusCode, StatusDescription(rw.statusCode))
	if _, err := rw.conn.Write([]byte(responseLine)); err != nil {
		return 0, fmt.Errorf("failed to write response line: %w", err)
	}

	// Get the header keys and sort them
	keys := make([]string, 0, len(rw.headers))
	for key := range rw.headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	rw.conn.SetWriteDeadline(time.Now().Add(rw.writeTimeout))
	// Write headers in sorted order
	for _, key := range keys {
		headerLine := fmt.Sprintf("%s: %s\r\n", key, rw.headers[key])
		if _, err := rw.conn.Write([]byte(headerLine)); err != nil {
			return 0, fmt.Errorf("failed to write header: %w", err)
		}
	}

	rw.conn.SetWriteDeadline(time.Now().Add(rw.writeTimeout))
	// Write \r\n between headers and body
	if _, err := rw.conn.Write([]byte("\r\n")); err != nil {
		return 0, fmt.Errorf("failed to write blank line after headers: %w", err)
	}

	rw.conn.SetWriteDeadline(time.Now().Add(rw.writeTimeout))
	// Write body
	n, err := rw.conn.Write(body)
	if err != nil {
		return n, fmt.Errorf("failed to write body: %w", err)
	}

	return n, nil
}
