package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type ViewDataResponse struct {
	Data   []Row        `json:"data"`
	Status []StatusItem `json:"status"`
}

type Row struct {
	ID              int64  `json:"id"`
	ProductID       string `json:"productID"`
	ProductName     string `json:"productName"`
	Amount          string `json:"amount"`
	CustomerName    string `json:"customerName"`
	Status          int    `json:"status"`
	TransactionDate string `json:"transactionDate"`
	CreateBy        string `json:"createBy"`
	CreateOn        string `json:"createOn"`
}

type StatusItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var viewData ViewDataResponse

func loadViewData() error {
	// Try to read from embedded data directory first (in binary)
	embedPath := filepath.Join("data", "viewData.json")
	data, err := os.ReadFile(embedPath)
	if err != nil {
		return fmt.Errorf("failed to read viewData.json: %w", err)
	}

	var parsed ViewDataResponse
	if err := json.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("invalid JSON in viewData.json: %w", err)
	}

	viewData = parsed
	return nil
}

func viewDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(viewData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func main() {
	if err := loadViewData(); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %v\n", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/view-data", viewDataHandler)
	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %v\n", err)
		os.Exit(1)
	}
}
