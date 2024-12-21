package http

import (
	"fmt"
	"io"
)

type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

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

// String converts the HTTPResponse to a string
func (r HTTPResponse) String() string {
	desc := StatusDescription(r.StatusCode)

	// Add status line
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.StatusCode, desc)

	// Add headers to response
	for k, v := range r.Headers {
		response += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	// Separate with headers, add body
	response += "\r\n"
	response += string(r.Body)

	return response
}

// Send response to connection
func (r HTTPResponse) Write(w io.Writer) error {
	_, err := w.Write([]byte(r.String()))
	return err
}
