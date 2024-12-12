package server

import (
	"fmt"
	"log"
	"net"
)

// TCP Server type
type Server struct {
	listenAddr string
	listener   net.Listener
	quitch     chan struct{}
	router     *HTTPRouter
}

// NewServer return new server object
func NewServer(listenAddr string, router *HTTPRouter) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		router:     router,
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

// handleConnections handles incoming connection (request)
func (s *Server) handleConnection(conn net.Conn) error {
	defer conn.Close()

	req, err := ParseRequest(conn)
	if err != nil {
		return fmt.Errorf("Failed to parse request: %s\n", err)
	}

	handler := s.router.GetHandler(req)

	var response HTTPResponse
	if handler == nil {
		response = HTTPResponse{
			StatusCode: 404,
		}
	} else {
		response = handler(req)
	}

	if err := response.Write(conn); err != nil {
		return fmt.Errorf("Error processing reponse: %s\n", err)
	}

	return nil
}
