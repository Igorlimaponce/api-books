package book

import (
	"strings"
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
}

type RequestCreateBook struct {
	Title       string
	Author      string
	Published   FlexibleDate
	Image       string
	Description string
}

type ResponseCreateBook struct {
	ID          uuid.UUID
	Title       string
	Author      string
	Published   FlexibleDate
	Image       string
	Description string
}

type RequestUpdateBook struct {
	ID          uuid.UUID
	Title       string
	Author      string
	Published   FlexibleDate
	Image       string
	Description string
}

// FlexibleDate é um tipo customizado que aceita múltiplos formatos de data
type FlexibleDate struct {
	time.Time
}

// UnmarshalJSON implementa json.Unmarshaler para aceitar diferentes formatos
func (fd *FlexibleDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	if s == "null" || s == "" {
		fd.Time = time.Time{}
		return nil
	}

	formats := []string{
		"2006-01-02",           // YYYY-MM-DD
		"02/01/2006",           // DD/MM/YYYY
		"01/02/2006",           // MM/DD/YYYY
		"2006-01-02T15:04:05Z", // ISO 8601
		time.RFC3339,           // RFC3339
	}

	var err error
	for _, format := range formats {
		fd.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return err
}
