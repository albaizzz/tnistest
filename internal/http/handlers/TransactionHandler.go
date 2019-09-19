package handlers

import (
	"net/http"

	"github.com/tnistest/internal/services"
)

type ITransactionHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
}

type TransactionHandler struct {
	TransactionSvc services.ITransactionHistorySvc
}

func NewTransactionHandler(transactionSvc services.ITransactionHistorySvc) *TransactionHandler {
	return &TransactionHandler{TransactionSvc: transactionSvc}
}

func (h *TransactionHandler) Get(w http.ResponseWriter, r *http.Request) {

}
