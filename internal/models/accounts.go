package models

import "time"

type Accounts struct {
	ID             int       `json:"id"`
	AccountNumber  string    `json:"account_number" valid:"required"`
	Name           string    `json:"name" valid:"required"`
	Address        string    `json:"address" valid:"required"`
	IdentityNumber string    `json:"identity_number" valid:"required"`
	IdentityType   uint32    `json:"identity_type" valid:"required"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"update_at"`
}
