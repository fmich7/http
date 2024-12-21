package http

import (
	"fmt"
	"log"
	"net"
)

// TCP Server
type Server struct {
	listenAddr string
	listener   net.Listener
	quitch     chan struct{}
	router     *HTTPRouter
}

// NewServer returns a new server object
func NewServer(listenAddr string, router *HTTPRouter) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		router:     router,
	}
}

// Start method starts the server
func (s *Server) Start() error {
	defer close(s.quitch)

	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("error starting tcp server: %s", err)
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

// handleConnections handles incoming connections (requests)
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := ParseRequest(conn)
	if err != nil {
		log.Printf("Failed to parse request from %s: %v", conn.RemoteAddr(), err)
		response := HTTPResponse{
			StatusCode: 400,
			Body:       []byte(StatusDescription(400) + "\n"),
		}
		response.Write(conn)
		return
	}

	handler, params := s.router.GetHandler(req)

	var response HTTPResponse
	if handler == nil {
		log.Printf("No handler found for path: %s", req.URL)
		response = HTTPResponse{
			StatusCode: 404,
			Body:       []byte(StatusDescription(404) + "\n"),
		}
	} else {
		response = handler(req, params)
	}

	if err := response.Write(conn); err != nil {
		log.Printf("Error writing response to %s: %v", conn.RemoteAddr(), err)
	}
}
