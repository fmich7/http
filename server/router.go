package server

type HTTPHandler func(HTTPRequest) HTTPResponse

type HTTPRouter struct {
	routes map[string]HTTPHandler
}

func NewHTTPRouter() *HTTPRouter {
	return &HTTPRouter{
		routes: make(map[string]HTTPHandler),
	}
}

func (s *HTTPRouter) AddEndpoint(method string, path string, handler HTTPHandler) {
	s.routes[method+path] = handler
}

func (s *HTTPRouter) GetHandler(req HTTPRequest) HTTPHandler {
	if handler, ok := s.routes[req.Method+req.URL]; ok {
		return handler
	}

	return nil
}
