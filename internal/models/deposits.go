package models

import "time"

type Deposits struct {
	AccountID         uint32    `json:"account_id" valid:"required"`
	Balance           float64   `json:"balance" valid:"required"`
	LastTransactionID int64     `json:"last_transaction_id" valid:"required"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"update_at"`
}
