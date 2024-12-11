package main

import (
	"fmt"
	"log"
	"net"
)

type HTTPRequest struct {
	Method string
	URL string
	ProtocolVersion string
	Headers map[string]string
	Body []byte
}

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

// TCP Server type
type Server struct {
	listenAddr string
	listener net.Listener
	quitch chan struct{}
}

// NewServer return new server object
func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch: make(chan struct{}),
	}
}

// Start method starts server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		log.Fatalln("Error starting tcp server:", err)
	}	
	defer listener.Close()
	
	fmt.Println("Running tcp server on address:", s.listenAddr)
	s.listener = listener

	go s.acceptLoop()

	<-s.quitch

	return nil
}

// acceptLoop method handles incoming connections
func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("New request from:", conn.RemoteAddr())
		go s.handleConnection(conn)
	}
}

// parseRequest gets infromations from incoming request
func parseRequest(conn net.Conn) (HTTPRequest, error) {
	for {
		req := HTTPRequest{
			Body: make([]byte, 1024),
		}
		_, err := conn.Read(req.Body)
		if err != nil {
			return HTTPRequest{}, fmt.Errorf("Read error: %s\n", err)
		}

		return req, nil
	}
}

// handleConnections handles incoming connection (request)
func (s *Server) handleConnection(conn net.Conn) error {
	defer conn.Close()

	req, err := parseRequest(conn)
	fmt.Println(string(req.Body))
	if err != nil {
		return fmt.Errorf("Failed to parse request: %s\n", err)
	}
	response := HTTPResponse{
		StatusCode: 200,	
	}	
	
	if err := response.Write(conn); err != nil {
		return fmt.Errorf("Error processing reponse: %s\n", err)
	}
	
	fmt.Println("Success")
	return nil
}

func main() {
	server := NewServer(":3000")
	server.Start()
}