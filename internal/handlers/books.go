package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

func BookById(mux chi.Router, sto *storage.BooksPostgresStorage) {
	mux.Get("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Printf("Invalid book ID: %v", err)
			return
		}
		book, err := sto.BookById(r.Context(), id)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error getting book: %v", err)
			return
		}

		bookJSON, err := json.Marshal(book)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error encoding book to JSON: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(bookJSON)
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
			log.Printf("Error adding book to storage: %v", err)
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
			log.Printf("Error encoding JSON: %v", err)
			return
		}
	})
}
func UpdateBook(mux chi.Router, sto *storage.BooksPostgresStorage) {
	mux.Put("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		var book model.Book
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Printf("Error converting ID to int: %v", err)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Printf("Error decoding JSON: %v", err)
			return
		}

		bookDb, err := sto.Update(r.Context(), book, id)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			log.Printf("Error updating book: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bookDb)
	})
}

func DeleteBook(mux chi.Router, sto *storage.BooksPostgresStorage) {
	mux.Delete("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		var idParam = chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Printf("Error converting id: %v", err)
			return
		}

		err = sto.Delete(r.Context(), id)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			log.Printf("Error deleting book: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
	})
}
