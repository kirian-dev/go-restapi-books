package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"restapi-books/internal/handlers"
	"restapi-books/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/matryer/is"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	dsn := "postgres://username:password@localhost:5432/testdb?sslmode=disable"
	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		panic("failed to connect to the database: " + err.Error())
	}

	_, err = db.Exec("CREATE TABLE books (id SERIAL PRIMARY KEY, title TEXT, published_at DATE)")
	if err != nil {
		panic("failed to create table: " + err.Error())
	}

	exitCode := m.Run()

	db.Close()

	os.Exit(exitCode)
}

func TestBooks(t *testing.T) {
	t.Run("returns 200", func(t *testing.T) {
		is := is.New(t)

		mux := chi.NewMux()

		storage := storage.NewBooksPostgresStorage(db)
		handlers.Books(mux, storage)
		code, _, body := makeGetRequest(mux, "/books")
		is.Equal(http.StatusOK, code)

		expected := `[]`
		is.Equal(expected, body)
	})
}

func makeGetRequest(handler http.Handler, target string) (int, http.Header, string) {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	result := res.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}

	return result.StatusCode, result.Header, string(bodyBytes)
}
