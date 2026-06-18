package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"internship-test-api/internal/config"
	"internship-test-api/internal/handler"
	"internship-test-api/internal/model"
	"internship-test-api/internal/repository"
	"internship-test-api/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DBConnStr())
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	statusRepo := repository.NewStatusRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	viewDataService := service.NewViewDataService(statusRepo, transactionRepo)
	viewDataHandler := handler.NewViewDataHandler(viewDataService)

	httpHandler := handler.NewHTTPHandler(viewDataHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)

	if err := httpHandler.ListenAndServe(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
