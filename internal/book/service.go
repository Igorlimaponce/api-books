package book

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type BookService struct {
	bookRepo *BookRepository
	timeout  time.Duration
	// Esse timeout vai ajudar na consulta ao banco (No caso definimos abaixo 2 segundos)
	// Ou seja o tempo maximo de um insert no banco por exemplo é de 2 sec, caso isso nao
	// ocorra, vai ser retornado um erro de estouro de tempo, isso é util para nao travar o programa e o banco
	// O travamento ocorreria por um consumo de recurso do banco e memoria desnecessario
}

func NewBookService(bookRepo *BookRepository) *BookService {
	return &BookService{
		bookRepo: bookRepo,
		timeout:  time.Duration(2) & time.Second,
	}
}

func (s *BookService) GetBookByID(ctx context.Context, id uuid.UUID) (*Book, error) {
	return s.bookRepo.GetBookByID(ctx, id)
}

func (s *BookService) DeleteBook(ctx context.Context, id uuid.UUID) error {
	return s.bookRepo.DeleteBook(ctx, id)
}

func (s *BookService) CreateBook(ctx context.Context, reqCreateBook RequestCreateBook) (*ResponseCreateBook, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	log.Printf("UserService.CreateBook - Starting book creation for: %s", reqCreateBook.Title)

	if reqCreateBook.Title == "" || reqCreateBook.Published.IsZero() || reqCreateBook.Author == "" {
		log.Printf("UserService.CreateBook - Validation failed: missing required fields")
		return nil, fmt.Errorf("title, published, and author are required")
	}

	b := &Book{
		Title:     reqCreateBook.Title,
		Published: reqCreateBook.Published,
		Author:    reqCreateBook.Author,
	}

	book, err := s.bookRepo.CreateBook(ctx, b)
	if err != nil {
		log.Printf("UserService.CreateBook - Error creating book: %v", err)
		return nil, fmt.Errorf("create book: %w", err)
	}
	log.Printf("UserService.CreateBook - Successfully created book with ID: %s", book.ID)
	return &ResponseCreateBook{
		ID:     book.ID,
		Title:  book.Title,
		Author: book.Author,
	}, nil
}

func (s *BookService) UpdateBook(ctx context.Context, book *Book) (*ResponseCreateBook, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	updatedBook, err := s.bookRepo.UpdateBook(ctx, book)
	if err != nil {
		return nil, fmt.Errorf("update book: %w", err)
	}
	return &ResponseCreateBook{
		ID:     updatedBook.ID,
		Title:  updatedBook.Title,
		Author: updatedBook.Author,
	}, nil
}
