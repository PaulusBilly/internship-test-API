package main

import (
	"log"

	"github.com/PaulusBilly/internship-test-API/internal/bootstrap"
	"github.com/PaulusBilly/internship-test-API/internal/config"
	"github.com/PaulusBilly/internship-test-API/internal/repository"
	"github.com/PaulusBilly/internship-test-API/internal/service"
	"github.com/PaulusBilly/internship-test-API/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBConnStr())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	repo := repository.NewPostgresRepository(db)
	svc := service.New(repo)

	// Seed data if needed
	if err := bootstrap.Seed(repo, "resources/viewData.json"); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	h := handler.New(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/view-data", h.GetViewData)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
