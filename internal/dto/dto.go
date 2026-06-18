package dto

type ViewDataResponse struct {
	Data   []Row        `json:"data"`
	Status []StatusItem `json:"status"`
}

type Row struct {
	ID               int64  `json:"id"`
	ProductID        string `json:"productID"`
	ProductName      string `json:"productName"`
	Amount           string `json:"amount"`
	CustomerName     string `json:"customerName"`
	Status           int    `json:"status"`
	TransactionDate  string `json:"transactionDate"`
	CreateBy         string `json:"createBy"`
	CreateOn         string `json:"createOn"`
}

type StatusItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
