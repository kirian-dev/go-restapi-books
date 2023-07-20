package server

import "restapi-books/internal/handlers"

func (s *Server) setupRoutes() {
	handlers.Books(s.mux)
}
