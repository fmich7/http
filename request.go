package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

// CRLF const
const CRLF = "\r\n"

// HTTPRequest represents the incoming HTTP request
type HTTPRequest struct {
	// Method of the request
	Method string

	// Target path of the request (e.g. "/api/something")
	URL string

	// HTTP version of the request
	ProtocolVersion string

	// Key-value pairs of HTTP headers included in the request
	Headers map[string]string

	// Body holds the raw content of the request body
	Body []byte

	// Params stores any route parameters extracted from the URL (e.g. "/users/{id}", gives {"id": "123"})
	Params map[string]string

	// Context for the request, which can carry deadlines
	ctx context.Context
}

// Context return's the request context
func (r *HTTPRequest) Context() context.Context {
	if r.ctx != nil {
		return r.ctx
	}
	return context.Background()
}

// WithContext returns copy of request with set ctx
func (r *HTTPRequest) WithContext(ctx context.Context) *HTTPRequest {
	if ctx == nil {
		panic("nil request's context")
	}
	n := new(HTTPRequest)
	*n = *r
	n.ctx = ctx
	return n
}

// ReadRequest reads request data bytes to buffer
func ReadRequest(conn net.Conn, readTimeout time.Duration) ([]byte, error) {
	// Set the read deadline
	conn.SetReadDeadline(time.Now().Add(readTimeout))
	defer conn.SetReadDeadline(time.Time{})

	// Read request
	var buf bytes.Buffer
	tmp := make([]byte, 1024)
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return nil, fmt.Errorf("read timeout occurred: %w", err)
			}

			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("read error: %w", err)
		}
		buf.Write(tmp[:n])
		// TODO: Improve this part
		if bytes.Contains(buf.Bytes(), []byte("\r\n\r\n")) {
			break
		}
	}

	return buf.Bytes(), nil
}

// ParseRequest gets infromations from incoming request
func ParseRequest(conn net.Conn, readTimeout time.Duration) (*HTTPRequest, error) {
	// Read data from connection
	reqData, err := ReadRequest(conn, readTimeout)
	if err != nil {
		return nil, err
	}

	// Split data and read method url protocol
	splitData := bytes.Split(reqData, []byte("\r\n"))
	reqLineVals := bytes.Split(splitData[0], []byte(" "))
	if len(reqLineVals) != 3 {
		return nil, errors.New("Invalid request line")
	}

	// Request attributes
	method := reqLineVals[0]
	url := reqLineVals[1]
	protocol := reqLineVals[2]
	headers := make(map[string]string)
	body := make([]byte, 0)
	// Get header values
	i := 1
	for i < len(splitData) {
		if len(splitData[i]) == 0 {
			i++
			break
		}

		headerLineValues := bytes.Split(splitData[i], []byte(": "))
		if len(headerLineValues) != 2 {
			return nil, errors.New("Invalid header entry")
		}

		headers[string(headerLineValues[0])] = string(headerLineValues[1])
		i++
	}

	// Copy data to body if there is any
	if i < len(splitData) {
		body = splitData[i]
	}

	return &HTTPRequest{
		Method:          string(method),
		URL:             string(url),
		ProtocolVersion: string(protocol),
		Headers:         headers,
		Body:            body,
	}, nil
}

// func (r HTTPRequest) String() string {
// 	headers := ""
// 	for k, v := range r.Headers {
// 		headers += fmt.Sprintf("%s: %s\n", k, v)
// 	}

// 	req := fmt.Sprintf("%s %s %s\n%s", r.Method, r.URL, r.ProtocolVersion, headers)

// 	if len(r.Body) > 0 {
// 		req += fmt.Sprintf("\n\n%s", r.Body)
// 	}

// 	return req
// }

// isEqualHTTPRequest compares that 2 requests are similar (used for testing)
func isEqualHTTPRequest(a, b HTTPRequest) bool {
	if a.Method != b.Method ||
		a.URL != b.URL ||
		a.ProtocolVersion != b.ProtocolVersion {
		return false
	}

	if len(a.Headers) != len(b.Headers) {
		return false
	}

	for k, v := range a.Headers {
		if vb, ok := b.Headers[k]; !ok || vb != v {
			return false
		}
	}

	if !bytes.Equal(a.Body, b.Body) {
		return false
	}

	return true
}
