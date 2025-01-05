package http

import "testing"

func TestNewHTTPRouter(t *testing.T) {
	var router *HTTPRouter = NewHTTPRouter()
	if len(router.routes) > 0 {
		t.Errorf("Routes length should be equal to 0 and not %d", len(router.routes))
	}
}

func TestHTTPRouter(t *testing.T) {
	// Initial setup
	router := NewHTTPRouter()
	dummyHandler := func(*HTTPRequest, ResponseWriter) {}

	// Register routes
	router.HandlerFunc("GET", "/hello", dummyHandler)
	router.HandlerFunc("GET", "/user/{id}", dummyHandler)
	router.HandlerFunc("POST", "/user/{id}/profile", dummyHandler)

	tests := []struct {
		name           string
		req            HTTPRequest
		expectHandler  bool
		expectedParams map[string]string
	}{
		{
			name: "Basic route match",
			req: HTTPRequest{
				Method: "GET",
				URL:    "/hello",
			},
			expectHandler:  true,
			expectedParams: map[string]string{},
		},
		{
			name: "Route with placeholder match",
			req: HTTPRequest{
				Method: "GET",
				URL:    "/user/123",
			},
			expectHandler: true,
			expectedParams: map[string]string{
				"id": "123",
			},
		},
		{
			name: "Route with multiple placeholders",
			req: HTTPRequest{
				Method: "POST",
				URL:    "/user/456/profile",
			},
			expectHandler: true,
			expectedParams: map[string]string{
				"id": "456",
			},
		},
		{
			name: "Method mismatch",
			req: HTTPRequest{
				Method: "POST",
				URL:    "/hello",
			},
			expectHandler: false,
		},
		{
			name: "Path length mismatch",
			req: HTTPRequest{
				Method: "GET",
				URL:    "/user/123/profile",
			},
			expectHandler: false,
		},
		{
			name: "No matching route",
			req: HTTPRequest{
				Method: "GET",
				URL:    "/nonexistent",
			},
			expectHandler: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := router.GetHandler(&tt.req)

			if tt.expectHandler && handler == nil {
				t.Errorf("Expected handler, but got nil")
			}

			if !tt.expectHandler && handler != nil {
				t.Errorf("Expected no handler, but got one")
			}

			if tt.expectHandler {
				for key, val := range tt.expectedParams {
					if tt.req.Params[key] != val {
						t.Errorf("Expected param %s to be %s, but got %s", key, val, tt.req.Params[key])
					}
				}
			}
		})
	}
}
