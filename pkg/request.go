package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
)

type HTTPRequest struct {
	Method          string
	URL             string
	ProtocolVersion string
	Headers         map[string]string
	Body            []byte
}

const CRLF = "\r\n"

// ReadRequest reads request data bytes to buffer
func ReadRequest(conn net.Conn) ([]byte, error) {
	var buf bytes.Buffer
	tmp := make([]byte, 1024)

	for {
		n, err := conn.Read(tmp)

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("read error: %w", err)
		}

		buf.Write(tmp[:n])
		if bytes.Contains(buf.Bytes(), []byte("\r\n\r\n")) {
			break
		}
	}
	return buf.Bytes(), nil
}

// ParseRequest gets infromations from incoming request
func ParseRequest(conn net.Conn) (HTTPRequest, error) {
	// Read data from connection
	reqData, err := ReadRequest(conn)
	if err != nil {
		return HTTPRequest{}, err
	}

	// Split data and read method url protocol
	splitData := bytes.Split(reqData, []byte("\r\n"))
	reqLineVals := bytes.Split(splitData[0], []byte(" "))
	if len(reqLineVals) != 3 {
		return HTTPRequest{}, errors.New("Invalid request line")
	}

	// Request attributes
	method := reqLineVals[0]
	url := reqLineVals[1]
	protocol := reqLineVals[2]
	headers := make(map[string]string)
	body := make([]byte, 0)

	// Get header values
	i := 1
	for ; i < len(splitData); i++ {
		if len(splitData[i]) == 0 {
			break
		}

		headerLineValues := bytes.Split(splitData[i], []byte(": "))
		if len(headerLineValues) != 2 {
			return HTTPRequest{}, errors.New("Invalid header entry")
		}

		headers[string(headerLineValues[0])] = string(headerLineValues[1])
	}

	// Copy data to body if there is any
	if i < len(splitData) {
		copy(body, splitData[i])
	}

	return HTTPRequest{
		Method:          string(method),
		URL:             string(url),
		ProtocolVersion: string(protocol),
		Headers:         headers,
		Body:            body,
	}, nil
}

func (r HTTPRequest) String() string {
	headers := ""
	for k, v := range r.Headers {
		headers += fmt.Sprintf("%s: %s\n", k, v)
	}

	req := fmt.Sprintf("%s %s %s\n%s", r.Method, r.URL, r.ProtocolVersion, headers)

	if len(r.Body) > 0 {
		req += fmt.Sprintf("\n\n%s", r.Body)
	}

	return req
}

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
