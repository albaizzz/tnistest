package handlers

import (
	"net/http"

	"github.com/tnistest/internal/services"
)

type IAccountHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
}

type AccountHandler struct {
	AccountSvc services.IAccountService
}

func NewAccountHandler(accountSvc services.IAccountService) *AccountHandler {
	return &AccountHandler{AccountSvc: accountSvc}
}

func (h *AccountHandler) Register(w http.ResponseWriter, r *http.Request) {

}
