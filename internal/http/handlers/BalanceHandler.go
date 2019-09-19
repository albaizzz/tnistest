package handlers

import (
	"net/http"

	"github.com/tnistest/internal/services"
)

type IBalanceHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
}

type BalanceHandler struct {
	BalanceSvc services.IBalanceService
}

func NewBalanceHandler(balanceSvc services.IBalanceService) *BalanceHandler {
	return &BalanceHandler{BalanceSvc: balanceSvc}
}

func (h *BalanceHandler) Get(w http.ResponseWriter, r *http.Request) {

}

func (h *BalanceHandler) Transact(w http.ResponseWriter, r *http.Request) {

}
