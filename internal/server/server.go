package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"books-api/internal/book"
	"books-api/internal/database"
)

type Server struct {
	port        int
	db          database.Service
	bookHandler *book.BookHandler
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// Inicializar database
	db := database.New()

	// Criar camadas seguindo a arquitetura (Repository -> Service -> Handler)
	bookRepo := book.NewBookRepository(db.GetDB()) // Você precisará adicionar um método GetDB() no service
	bookService := book.NewBookService(bookRepo)
	bookHandler := book.NewBookHandler(bookService)

	NewServer := &Server{
		port:        port,
		db:          db,
		bookHandler: bookHandler,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
