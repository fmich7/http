package http

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Server represents an HTTP server that listens and handles requests
type Server struct {
	listenAddr   string
	listener     net.Listener
	quitch       chan struct{}
	startch      chan struct{}
	running      bool
	runningMu    sync.Mutex
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
		startch:      make(chan struct{}),
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
	close(s.startch)
	s.setRunning(true)

	// Wait for a SIGINT or SIGTERM signal to gracefully shut down the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	s.Stop()
	log.Println("Server stopped")

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
		s.setRunning(false)
	case <-time.After(time.Second):
		log.Println("Timed out waiting for connections to finish.")
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

// setRunning method sets server running state
func (s *Server) setRunning(state bool) {
	s.runningMu.Lock()
	defer s.runningMu.Unlock()
	s.running = state
}

// setRunning method gets server running state
func (s *Server) isRunning() bool {
	s.runningMu.Lock()
	defer s.runningMu.Unlock()
	return s.running
}
