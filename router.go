package http

import "strings"

// HTTPHandler defines signature of func that handles requests
type HTTPHandler func(*HTTPRequest, ResponseWriter)

// Route represents a single HTTP route in the router
type Route struct {
	method     string
	pathParts  []string
	handler    HTTPHandler
	paramNames []string
}

// HTTPRouter is a router for managing routes and their handlers
type HTTPRouter struct {
	routes []Route
}

// NewHTTPRouter return new HTTPRouter
func NewHTTPRouter() *HTTPRouter {
	return &HTTPRouter{
		routes: make([]Route, 0),
	}
}

// HandlerFunc adds a new route to the HTTPRouter
// - method: HTTP method for the route "GET", "POST", etc.
// - path: URL path for the route (dynamic parameters "/users/{id}")
// - handler: Function to handle requests
func (s *HTTPRouter) HandlerFunc(method string, path string, handler HTTPHandler) {
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	paramNames := []string{}

	for _, part := range pathParts {
		if part[0] == '{' && part[len(part)-1] == '}' {
			paramNames = append(paramNames, part[1:len(part)-1])
		} else {
			paramNames = append(paramNames, "")
		}
	}

	s.routes = append(s.routes, Route{
		method:     method,
		pathParts:  pathParts,
		handler:    handler,
		paramNames: paramNames,
	})
}

// GetHandler returns the HTTP handler that is appropiate for given request
func (s *HTTPRouter) GetHandler(req *HTTPRequest) HTTPHandler {
	reqParts := strings.Split(strings.Trim(req.URL, "/"), "/")

	for _, route := range s.routes {
		if req.Method != route.method {
			continue
		}

		if len(reqParts) != len(route.pathParts) {
			continue
		}

		params := make(map[string]string)
		matches := true

		for i, routePart := range route.pathParts {
			if route.paramNames[i] != "" {
				params[route.paramNames[i]] = reqParts[i]
			} else if reqParts[i] != routePart {
				matches = false
				break
			}
		}

		if matches {
			req.Params = params
			return route.handler
		}
	}

	return nil
}
