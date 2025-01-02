package http

import (
	"context"
	"log"
	"time"
)

// TimeoutHandler timeouts request that take too long to finish
func TimeoutHandler(org HTTPHandler, dt time.Duration) HTTPHandler {
	return func(r *HTTPRequest, w ResponseWriter) {
		ctx, cancel := context.WithTimeout(r.Context(), dt)
		defer cancel()

		// Create a request with the new context
		r = r.WithContext(ctx)

		done := make(chan struct{})

		go func() {
			defer close(done)
			org(r, w)
		}()

		select {
		case <-done:
			// Handler completed within the timeout
			return
		case <-ctx.Done():
			// Timeout or cancellation occurred
			log.Printf("Timeout occurred for request to %s", r.URL)
			w.SetStatus(504)
			w.Write([]byte(StatusDescription(504)))
		}
	}
}
