package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// ViewDataResponse matches the JSON schema of the original viewData.json
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

// embeddedViewData is the exact content of src/main/resources/viewData.json
var embeddedViewData = `{
  "data": [
    {
      "id": 1372,
      "productID": "10001",
      "productName": "Test 1",
      "amount": "1000",
      "customerName": "abc",
      "status": 0,
      "transactionDate": "2022-07-10 11:14:52",
      "createBy": "abc",
      "createOn": "2022-07-10 11:14:52"
    },
    {
      "id": 1373,
      "productID": "10002",
      "productName": "Test 2",
      "amount": "2000",
      "customerName": "abc",
      "status": 0,
      "transactionDate": "2022-07-11 13:14:52",
      "createBy": "abc",
      "createOn": "2022-07-10 13:14:52"
    },
    {
      "id": 1374,
      "productID": "10001",
      "productName": "Test 1",
      "amount": "1000",
      "customerName": "abc",
      "status": 0,
      "transactionDate": "2022-08-10 12:14:52",
      "createBy": "abc",
      "createOn": "2022-07-10 12:14:52"
    },
    {
      "id": 1375,
      "productID": "10002",
      "productName": "Test 2",
      "amount": "1000",
      "customerName": "abc",
      "status": 1,
      "transactionDate": "2022-08-10 13:10:52",
      "createBy": "abc",
      "createOn": "2022-07-10 13:10:52"
    },
    {
      "id": 1376,
      "productID": "10001",
      "productName": "Test 1",
      "amount": "1000",
      "customerName": "abc",
      "status": 0,
      "transactionDate": "2022-08-10 13:11:52",
      "createBy": "abc",
      "createOn": "2022-07-10 13:11:52"
    },
    {
      "id": 1377,
      "productID": "10002",
      "productName": "Test 2",
      "amount": "2000",
      "customerName": "abc",
      "status": 0,
      "transactionDate": "2022-08-12 13:14:52",
      "createBy": "abc",
      "createOn": "2022-07-10 13:14:52"
    },
    {
      "id": 1378,
      "productID": "10001",
      "productName": "Test 1",
      "amount": "1000",
      "customerName": "abc",
      "status": 0,
      "transactionDate": "2022-08-12 14:11:52",
      "createBy": "abc",
      "createOn": "2022-07-10 14:11:52"
    },
    {
      "id": 1379,
      "productID": "10002",
      "productName": "Test 2",
      "amount": "1000",
      "customerName": "abc",
      "status": 1,
      "transactionDate": "2022-09-13 11:14:52",
      "createBy": "abc",
      "createOn": "2022-07-10 11:14:52"
    },
    {
      "id": 1380,
      "productID": "10001",
      "productName": "Test 1",
      "amount": "1000",
      "customerName": "abc",
      "status": 0,
      "transactionDate": "2022-09-13 13:14:52",
      "createBy": "abc",
      "createOn": "2022-07-10 13:14:52"
    },
    {
      "id": 1381,
      "productID": "10002",
      "productName": "Test 2",
      "amount": "2000",
      "customerName": "abc",
      "status": 0,
      "transactionDate": "2022-09-14 09:11:52",
      "createBy": "abc",
      "createOn": "2022-07-10 09:11:52"
    },
    {
      "id": 1382,
      "productID": "10001",
      "productName": "Test 1",
      "amount": "1000",
      "customerName": "abc",
      "status": 0,
      "transactionDate": "2022-09-14 10:14:52",
      "createBy": "abc",
      "createOn": "2022-07-10 10:14:52"
    },
    {
      "id": 1383,
      "productID": "10002",
      "productName": "Test 2",
      "amount": "1000",
      "customerName": "abc",
      "status": 1,
      "transactionDate": "2022-08-15 13:14:52",
      "createBy": "abc",
      "createOn": "2022-07-10 13:14:52"
    }
  ],
  "status": [
    {
      "id": 0,
      "name": "SUCCESS"
    },
    {
      "id": 1,
      "name": "FAILED"
    }
  ]
}`

var viewData ViewDataResponse

func init() {
	if err := json.Unmarshal([]byte(embeddedViewData), &viewData); err != nil {
		log.Fatalf("failed to parse embedded viewData.json: %v", err)
	}
}

func viewDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(viewData); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/api/view-data", viewDataHandler)

	port := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}

	fmt.Printf("Starting server on port 8080...\n")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("server failed to start on port 8080: %v", err)
	}
}
