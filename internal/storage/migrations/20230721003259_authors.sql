-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS authors (
  ID SERIAL PRIMARY KEY,
  NAME VARCHAR(255) NOT NULL,
	books_id INTEGER REFERENCES books (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS authors;
-- +goose StatementEnd
