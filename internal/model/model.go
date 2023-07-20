package mop

import "time"

type Book struct {
	ID          int64
	Title       string
	PublishedAt time.Time
}

type Author struct {
	ID     int64
	BookID int64
	Name   string
}
