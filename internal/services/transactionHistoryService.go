package services

import (
	"database/sql"

	"github.com/tnistest/internal/models"
	"github.com/tnistest/internal/repositories"
)

type ITransactionHistorySvc interface {
	GetByAccountNumber(accountNumber string) models.APIResponse
}

type TransactionHistorySvc struct {
	TransactionHistRepo repositories.ITransactionRepository
}

func NewTransactionHistory(transactionRepo repositories.ITransactionRepository) *TransactionHistorySvc {
	return &TransactionHistorySvc{TransactionHistRepo: transactionRepo}
}

func (t *TransactionHistorySvc) GetByAccountNumber(accountNumber string) models.APIResponse {
	db, err := b.DB.GetDB()
	if err != nil {
		return
	}

	rows, err := db.Query(`select d.balance, a.account_number from deposit d inner join accounts a on d.account_id = a.id
		where d.account_id = ?`, accountID)

	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(
			&balance.Balance,
			&balance.AccountNumber,
		)
	}
	if err != nil && err != sql.ErrNoRows {
		return
	}
	return balance, nil
}
