package services

import (
	"github.com/tnistest/internal/consts"
	"github.com/tnistest/internal/models"
	"github.com/tnistest/internal/repositories"
	"github.com/tnistest/pkg/helpers"
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
	transList, err := t.TransactionHistRepo.GetByAccount(accountNumber)
	if err != nil {
		return helpers.GetAPIResponse(consts.APIErrorUnknown, nil)
	}
	return helpers.GetAPIResponse(consts.APIGeneralSuccess, transList)
}
