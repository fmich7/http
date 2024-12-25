package http

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// TCP Server
type Server struct {
	listenAddr   string
	listener     net.Listener
	quitch       chan struct{}
	router       *HTTPRouter
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewServer returns a new server object
func NewServer(listenAddr string, router *HTTPRouter) *Server {
	return &Server{
		listenAddr:   listenAddr,
		quitch:       make(chan struct{}),
		router:       router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
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
		log.Println("New request from:", conn.RemoteAddr())
		go s.handleConnection(conn)
	}
}

// handleConnections handles incoming connections (requests)
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := ParseRequest(conn, s.ReadTimeout)
	rw := NewResponseWriter(conn, s.WriteTimeout)
	if err != nil {
		// Send 408 if read took too long
		if strings.Contains(err.Error(), "read timeout") {
			log.Printf("Read timeout %s: %v", conn.RemoteAddr(), err)
			rw.SetStatus(408)
			rw.Write([]byte(StatusDescription(408) + "\n"))
		}
		log.Printf("Failed to parse request from %s: %v", conn.RemoteAddr(), err)
		rw.SetStatus(400)
		rw.Write([]byte(StatusDescription(400) + "\n"))
		return
	}

	handler, params := s.router.GetHandler(req)

	if handler == nil {
		log.Printf("No handler found for path: %s", req.URL)
		rw.SetStatus(404)
		rw.Write([]byte(StatusDescription(404) + "\n"))
	} else {
		handler(req, rw, params)
	}
}
