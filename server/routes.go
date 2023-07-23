package server

import (
	"restapi-books/internal/handlers"
)

func (s *Server) setupRoutes() {
	handlers.Books(s.mux, s.storage)
	handlers.BookById(s.mux, s.storage)
	handlers.AddBook(s.mux, s.storage)
	handlers.UpdateBook(s.mux, s.storage)
	handlers.DeleteBook(s.mux, s.storage)
}
