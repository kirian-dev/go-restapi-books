package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Books(mux chi.Router) {
	mux.Get("/books", func(w http.ResponseWriter, r *http.Request) {

	})
}
