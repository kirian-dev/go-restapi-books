package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"restapi-books/internal/model"
	"restapi-books/internal/storage"

	"github.com/go-chi/chi/v5"
)

func Books(mux chi.Router, sto *storage.BooksPostgresStorage) {
	mux.Get("/books", func(w http.ResponseWriter, r *http.Request) {
		books, err := sto.Books(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error getting books from storage: %v", err)
			return
		}
		booksJSON, err := json.Marshal(books)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error marshalling books to JSON: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(booksJSON)
	})
}

func AddBook(mux chi.Router, sto *storage.BooksPostgresStorage) {
	mux.Post("/books", func(w http.ResponseWriter, r *http.Request) {
		var book model.Book
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Printf("Error decoding JSON: %v", err)
			return
		}

		id, err := sto.Add(r.Context(), book)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			log.Printf("error adding book to storage: %v", err)
			return
		}

		response := struct {
			ID int64 `json:"id"`
		}{
			ID: *id,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			log.Printf("error encoding json: %v", err)
			return
		}
	})
}
