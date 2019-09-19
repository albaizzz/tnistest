package models

type Balance struct {
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
}
