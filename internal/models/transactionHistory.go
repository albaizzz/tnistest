package models

import "time"

type TransactionHistory struct {
	ID                int       `json:"id"`
	TransactionAmount float64   `json:"transaction_amount" valid:"required"`
	FinalBalance      float64   `json:"final_balance"`
	TransactionType   uint      `json:"transaction_type" valid:"required"`
	OfficerID         uint      `json:"office_id" valid:"required"`
	AccountID         uint32    `json:"account_id" valid:"required"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"update_at"`
}
