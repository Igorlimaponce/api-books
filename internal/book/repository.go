package book

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type BookRepositoryInterface interface {
	GetBookByID(ctx context.Context, id uuid.UUID) (*Book, error)
	CreateBook(ctx context.Context, book *Book) (*Book, error)
	UpdateBook(ctx context.Context, book *Book) (*Book, error)
	DeleteBook(ctx context.Context, id uuid.UUID) error
}

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetBookByID(ctx context.Context, id uuid.UUID) (*Book, error) {
	query := `
		SELECT id, title, author, published, image, description, created_at, updated_at, deleted_at
		FROM books
		WHERE id = $1
	`

	var book Book
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Published,
		&book.Image,
		&book.Description,
		&book.Created_at,
		&book.Updated_at,
		&book.Deleted_at,
	)

	// Acima nós mandamos todos os valores retornados para a struc book,
	// Fazemos tudo com & na frente para passar o endereco e nao a struct por completa
	// Isso é um recurso de Go que ajuda a melhorar a performance do programa

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query book by id: %w", err)
	}

	return &book, nil
}

func (r *BookRepository) CreateBook(ctx context.Context, book *Book) (*Book, error) {
	query := `
		INSERT INTO books (title, author, published_date) 
		VALUES (&1, &2, &3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		book.Title, book.Author, book.Published).Scan(&book.ID, &book.Created_at, &book.Updated_at)

	if err != nil {
		return nil, fmt.Errorf("insert book: %w", err)
	}
	return book, nil
}

func (r *BookRepository) UpdateBook(ctx context.Context, book *Book) (*Book, error) {
	query := `
		UPDATE books
		SET title = &1, author = &2, published_date = &3,
		image_url = &4, description = &5, updated_at = NOW()
		WHERE id = &6
		RETURNING id, title, author, published_date, image_url, description, created_at, updated_at, deleted_at
	`
	var updatedBook Book
	err := r.db.QueryRowContext(
		ctx, query,
		book.Title, book.Author, book.Published, book.Image, book.Description, book.ID).Scan(
		&updatedBook.ID,
		&updatedBook.Title,
		&updatedBook.Author,
		&updatedBook.Published,
		&updatedBook.Image,
		&updatedBook.Description,
		&updatedBook.Created_at,
		&updatedBook.Updated_at,
		&updatedBook.Deleted_at,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("book not found")
		}
		return nil, fmt.Errorf("update book: %w", err)
	}

	return &updatedBook, nil
}

func (r *BookRepository) DeleteBook(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM books
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}
