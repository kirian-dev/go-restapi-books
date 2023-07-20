package server

import "restapi-books/handlers"

func (s *Server) setupRoutes() {
	handlers.Books(s.mux)
}
