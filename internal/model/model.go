package model

import "time"

type Book struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	PublishedAt time.Time `json:"published_at"`
}

type Author struct {
	ID     int64
	BookID int64
	Name   string
}
