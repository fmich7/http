package http

import (
	"testing"
	"time"
)

// startTestServer starts a test server and returns it along with the port it is running on
func startTestServer(t *testing.T, router *HTTPRouter) (*Server, int) {
	s := NewServer(":0", router)

	go func() {
		if err := s.Start(); err != nil {
			t.Errorf("Failed to start server: %v", err)
		}
	}()

	select {
	case <-s.startch:
		return s, s.GetPort()
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Server did not start")
		return nil, 0
	}
}
