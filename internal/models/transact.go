package models

type Transac struct {
	AccountID       uint32  `json:"account_id" valid:"required"`
	OfficerID       uint    `json:"officer_id" valid:"required"`
	Amount          float64 `json:"amount" valid:"required"`
	TransactionType uint    `json:"transaction_type" valid:"required"`
}
