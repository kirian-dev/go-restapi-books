package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"restapi-books/internal/storage"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Server struct {
	address string
	mux     chi.Router
	log     *zap.Logger
	server  *http.Server
	storage *storage.BooksPostgresStorage
}

type Options struct {
	Host string
	Log  *zap.Logger
	Port int
}

func New(opts Options, storage *storage.BooksPostgresStorage) *Server {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}
	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := chi.NewMux()

	return &Server{
		address: address,
		mux:     mux,
		log:     opts.Log,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
		storage: storage,
	}
}

func (s *Server) Start() error {
	s.setupRoutes()

	s.log.Info("Starting", zap.String("address", s.address))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error starting server: %v", err)

	}
	return nil
}

func (s *Server) Stop() error {
	s.log.Info("Stopping", zap.String("address", s.address))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %v", err)
	}
	return nil
}
