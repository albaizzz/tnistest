package models

type UserTransactionHistory struct {
	TransactionAmount float64 `json:"transaction_amount" valid:"required"`
	FinalBalance      float64 `json:"final_balance"`
	TransactionType   uint    `json:"transaction_type" valid:"required"`
	AccountNumber     string  `json:"account_number"`
}
