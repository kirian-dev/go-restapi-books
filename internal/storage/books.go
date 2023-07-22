package storage

import (
	"context"
	"fmt"
	"restapi-books/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type BooksPostgresStorage struct {
	db *sqlx.DB
}

type dbBook struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	PublishedAt time.Time `db:"published_date"`
}

func NewBooksPostgresStorage(db *sqlx.DB) *BooksPostgresStorage {
	return &BooksPostgresStorage{
		db: db,
	}
}

func (s *BooksPostgresStorage) Books(ctx context.Context) ([]model.Book, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connect to db: %w", err)
	}

	defer conn.Close()

	var books []dbBook
	if err := conn.SelectContext(ctx, &books, "SELECT * FROM books"); err != nil {
		return nil, fmt.Errorf("error select all books: %v", err)
	}

	var modelBooks []model.Book
	for _, dbb := range books {
		book := model.Book{
			ID:          dbb.ID,
			Title:       dbb.Title,
			PublishedAt: dbb.PublishedAt,
		}
		modelBooks = append(modelBooks, book)
	}
	return modelBooks, nil
}

func (s *BooksPostgresStorage) BookById(ctx context.Context, id int) (*model.Book, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connect to db: %w", err)
	}

	defer conn.Close()

	var book dbBook

	if err := conn.SelectContext(ctx, &book, `SELECT * FROM book WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("error fetching book from db: %w", err)
	}

	return (*model.Book)(&book), nil
}

func (s *BooksPostgresStorage) Add(ctx context.Context, book model.Book) (*int64, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connect to db: %w", err)
	}

	defer conn.Close()

	var id int64

	row := conn.QueryRowxContext(ctx, `INSERT INTO books (title, published_at) VALUES ($1, $2) RETURNING id`, book.Title, book.PublishedAt)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("error inserting book into db: %w", err)
	}

	if err := row.Scan(&id); err != nil {
		return nil, fmt.Errorf("error inserting book into db: %w", err)
	}

	return &id, nil

}

func (s *BooksPostgresStorage) Update(ctx context.Context, book model.Book, id int64) (*model.Book, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error updating book from db: %w", err)
	}

	defer conn.Close()

	var dbBook dbBook

	if err := conn.SelectContext(ctx, &dbBook, `SELECT * FROM books WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("error fetching book from db: %w", err)
	}

	dbBook.Title = book.Title
	dbBook.PublishedAt = book.PublishedAt

	row := conn.QueryRowxContext(ctx, `UPDATE books SET title = $1, PUBLISHED_DATE = $2 WHERE id = $3`, dbBook.Title, dbBook.PublishedAt, id)

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("error updating book from db: %w", err)
	}

	return &model.Book{
		ID:          dbBook.ID,
		Title:       dbBook.Title,
		PublishedAt: dbBook.PublishedAt,
	}, nil
}
