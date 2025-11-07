package book

import (
	"books-api/util"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type BookHandler struct {
	bookService *BookService
}

func NewBookHandler(bookService *BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req RequestCreateBook

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateBook - JSON decode error: %v", err)
		util.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	log.Printf("CreateBook - Request received: title=%s, author=%s", req.Title, req.Author)

	book, err := h.bookService.CreateBook(r.Context(), req)
	if err != nil {
		log.Printf("CreateBook - Service error: %v", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to create book")
		return
	}

	log.Printf("CreateBook - Success: book created with ID=%s,Title=%s,Author=%s", book.ID, book.Title, book.Author)

	util.WriteJSON(w, http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Printf("UpdateBook - JSON decode error: %v", err)
		util.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	log.Printf("UpdateBook - Request received: title=%s, author=%s", book.Title, book.Author)

	updatedBook, err := h.bookService.UpdateBook(r.Context(), &book)
	if err != nil {
		log.Printf("UpdateBook - Service error: %v", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to update book")
		return
	}

	log.Printf("UpdateBook - Success: book updated with ID=%s, Title=%s, Author=%s", updatedBook.ID, updatedBook.Title, updatedBook.Author)

	util.WriteJSON(w, http.StatusOK, updatedBook)
}

func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	var idParam = r.URL.Query().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("GetBookByID - Invalid UUID: %v", err)
		util.WriteError(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	log.Printf("GetBookByID - Request received: ID=%s", id)

	book, err := h.bookService.GetBookByID(r.Context(), id)
	if err != nil {
		log.Printf("GetBookByID - Service error: %v", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to get book")
		return
	}

	log.Printf("GetBookByID - Success: book retrieved with ID=%s, Title=%s, Author=%s", book.ID, book.Title, book.Author)

	util.WriteJSON(w, http.StatusOK, book)

}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	var idParam = r.URL.Query().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("DeleteBook - Invalid UUID: %v", err)
		util.WriteError(w, http.StatusBadRequest, "invalid book ID")
		return
	}
	log.Printf("DeleteBook - Request received: ID=%s", id)

	if err := h.bookService.DeleteBook(r.Context(), id); err != nil {
		log.Printf("DeleteBook - Service error: %v", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to delete book")
		return
	}
	log.Printf("DeleteBook - Success: book deleted with ID=%s", id)
	util.WriteJSON(w, http.StatusNoContent, nil)
}
