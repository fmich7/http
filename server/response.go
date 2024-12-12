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

func StatusDescription(code int) (string, error) {
	m := map[int]string{
		200: "OK",
		404: "Not Found",
	}

	if desc, ok := m[code]; ok {
		return desc, nil
	}

	return "", fmt.Errorf("Invalid response code (%d)", code)
}

// Convert HTTPResponse object to string
func (r HTTPResponse) String() string {
	desc, err := StatusDescription(r.StatusCode)
	if err != nil {
		panic(err)
	}
	// Add status line
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n\r\n", r.StatusCode, desc)

	// Add headers to response
	for k, v := range r.Headers {
		response += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	// Add data to response body
	response += "\r\n" + string(r.Body)

	return response
}

// Send response to connection
func (r HTTPResponse) Write(conn net.Conn) error {
	_, err := conn.Write([]byte(r.String()))
	return err
}
