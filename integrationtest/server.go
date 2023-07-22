package integrationtest

import (
	"net/http"
	"testing"
	"time"

	"restapi-books/internal/storage"
	"restapi-books/pkg/db"
	"restapi-books/server"
)

func CreateServer() func() {
	db, err := db.CreateDbConnection()
	if err != nil {
		return nil
	}
	storage := storage.NewBooksPostgresStorage(db)
	s := server.New(server.Options{
		Host: "localhost",
		Port: 8081,
	}, storage)

	go func() {
		if err := s.Start(); err != nil {
			panic(err)
		}
	}()

	for {
		_, err := http.Get("http://localhost:8081/")
		if err != nil {
			break
		}
		time.Sleep(5 * time.Microsecond)
	}

	return func() {
		if err := s.Stop(); err != nil {
			panic(err)
		}
	}
}

func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
}
