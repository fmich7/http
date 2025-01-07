package http

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Server represents an HTTP server that listens and handles requests
type Server struct {
	listenAddr   string
	listener     net.Listener
	quitch       chan struct{}
	wg           sync.WaitGroup
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
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("error starting tcp server: %s", err)
	}

	fmt.Println("Running tcp server on address:", s.listenAddr)
	s.listener = listener

	s.wg.Add(1)
	go s.acceptLoop()

	return nil
}

// Stop method stops the server
func (s *Server) Stop() {
	close(s.quitch)
	s.listener.Close()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		log.Println("Timed out waiting for connections to finish.")
		return
	}
}

// acceptLoop method accepts incoming requests
func (s *Server) acceptLoop() {
	defer s.wg.Done()

	for {
		select {
		case <-s.quitch:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}
			log.Println("New request from:", conn.RemoteAddr())
			go s.handleConnection(conn)
		}
	}
}

// GetPort method return a port of running server
func (s *Server) GetPort() int {
	return s.listener.Addr().(*net.TCPAddr).Port
}

// handleConnections handles requests
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
			return
		}
		log.Printf("Failed to parse request from %s: %v", conn.RemoteAddr(), err)
		rw.SetStatus(400)
		rw.Write([]byte(StatusDescription(400) + "\n"))
		return
	}

	handler := s.router.GetHandler(req)

	if handler == nil {
		log.Printf("No handler found for path: %s", req.URL)
		rw.SetStatus(404)
		rw.Write([]byte(StatusDescription(404) + "\n"))
	} else {
		handler(req, rw)
	}
}
