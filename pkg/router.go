package http

import "strings"

// request, params -> response
type HTTPHandler func(HTTPRequest, ResponseWriter, map[string]string)

type Route struct {
	method     string
	pathParts  []string
	handler    HTTPHandler
	paramNames []string
}

type HTTPRouter struct {
	routes []Route
}

// NewHTTPRouter return new HTTPRouter
func NewHTTPRouter() *HTTPRouter {
	return &HTTPRouter{
		routes: make([]Route, 0),
	}
}

// HandlerFunc binds handler to route
// Use placeholders such as {variable}
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

// GetHandler returns the handler that matches url and extracted variables
func (s *HTTPRouter) GetHandler(req HTTPRequest) (HTTPHandler, map[string]string) {
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
			return route.handler, params
		}
	}

	return nil, nil
}
