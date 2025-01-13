package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// Server represents an HTTP server that listens and handles requests
type Server struct {
	listenAddr      string
	listener        net.Listener
	startch         chan struct{}
	running         int32
	wg              sync.WaitGroup
	router          *HTTPRouter
	serverCtx       context.Context
	cancelFunc      context.CancelFunc
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// NewServer returns a new server object
func NewServer(listenAddr string, router *HTTPRouter) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		listenAddr:      listenAddr,
		startch:         make(chan struct{}),
		router:          router,
		serverCtx:       ctx,
		cancelFunc:      cancel,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
		ShutdownTimeout: 10 * time.Second,
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
	s.Shutdown()
	log.Println("Server Shutdownped")

	return nil
}

// Shutdown method Shutdowns the server
func (s *Server) Shutdown() {
	s.cancelFunc()
	s.listener.Close()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.setRunning(false)
	case <-time.After(s.ShutdownTimeout):
		log.Println("Timed out waiting for connections to finish")
	}
}

// acceptLoop method accepts incoming requests
func (s *Server) acceptLoop() {
	defer s.wg.Done()

	for {
		select {
		case <-s.serverCtx.Done():
			log.Println("Closing accepting loop")
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}
			log.Println("New request from:", conn.RemoteAddr())
			s.wg.Add(1)
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
	defer s.wg.Done()

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

	req.ctx = s.serverCtx
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
	if state {
		atomic.StoreInt32(&s.running, 1)
	} else {
		atomic.StoreInt32(&s.running, 0)
	}
}

// setRunning method gets server running state
func (s *Server) isRunning() bool {
	return atomic.LoadInt32(&s.running) == 1
}
