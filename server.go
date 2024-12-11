package main

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
}

// NewServer return new server object
func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
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
