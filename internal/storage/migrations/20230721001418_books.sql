-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS books (
  ID SERIAL PRIMARY KEY,
  TITLE VARCHAR(255) NOT NULL,
  PUBLISHED_AT TIMESTAMP NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS books;
-- +goose StatementEnd
