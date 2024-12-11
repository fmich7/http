package main

import (
	"fmt"
	"net"
)

type HTTPResponse struct {
	StatusCode int
	Headers map[string]string
	Body []byte
}

// Convert HTTPResponse object to string
func (r HTTPResponse) String() string {
	// Add status line
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n\r\n", r.StatusCode, "OK")
	
	// Add headers to response
	for k, v := range r.Headers {
		response += fmt.Sprintf("%s: %s\r\n",k, v)
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
