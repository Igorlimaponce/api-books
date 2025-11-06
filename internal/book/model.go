package book

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID
	Title       string
	Author      string
	Published   time.Time
	Image       string
	Description string
	Created_at  time.Time
	Updated_at  time.Time
	Deleted_at  time.Time
}

type RequestCreateBook struct {
	Title     string
	Author    string
	Published time.Time
}

type ResponseCreateBook struct {
	ID     uuid.UUID
	Title  string
	Author string
}
