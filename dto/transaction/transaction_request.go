package transactiondto

type TransactionRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Total   int  `json:"total"`
}
