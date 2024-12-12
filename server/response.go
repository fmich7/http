package server

import (
	"fmt"
	"net"
)

type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// StatusDescription returns a status description for the given status code
func StatusDescription(code int) (string, error) {
	switch code {
	case 200:
		return "OK", nil
	case 404:
		return "Not Found", nil
	case 500:
		return "Internal Server Error", nil
	default:
		return "", fmt.Errorf("Unknown status code: %d", code)
	}
}

// String converts the HTTPResponse to a string
func (r HTTPResponse) String() string {
	desc, err := StatusDescription(r.StatusCode)
	if err != nil {
		panic(err)
	}

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
func (r HTTPResponse) Write(conn net.Conn) error {
	_, err := conn.Write([]byte(r.String()))
	return err
}
