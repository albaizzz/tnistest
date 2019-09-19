package services

import (
	"fmt"
	"log"

	"github.com/tnistest/internal/consts"
	"github.com/tnistest/internal/models"
	"github.com/tnistest/internal/repositories"
	"github.com/tnistest/pkg/helpers"
)

type IBalanceService interface {
	Get(accountNumber string) (resp models.APIResponse)
	Transact(transact models.Transac) (resp models.APIResponse)
}
type BalanceService struct {
	BalanceRepository repositories.IBalanceRepository
	TransactionHist   repositories.ITransactionRepository
}

func NewBalanceService(balanceRepo repositories.IBalanceRepository, transactionHistRepo repositories.ITransactionRepository) *BalanceService {
	return &BalanceService{BalanceRepository: balanceRepo, TransactionHist: transactionHistRepo}
}

func (t *BalanceService) Get(accountNumber string) (resp models.APIResponse) {
	balance, err := t.BalanceRepository.GetByAccountNumber(accountNumber)
	if err != nil {
		//write to log
		resp = helpers.GetAPIResponse(consts.APIErrorUnknown, nil)
		return
	}
	resp = helpers.GetAPIResponse(consts.APIGeneralSuccess, balance)
	return
}

func (t *BalanceService) Transact(transact models.Transac) (resp models.APIResponse) {
	//get balance
	balance, err := t.BalanceRepository.GetById(transact.AccountID)
	if err != nil {
		resp = helpers.GetAPIResponse(consts.APIErrorUnknown, nil)
		return
	}

	switch transact.TransactionType {
	case consts.CashOut:
		transactionStatus, deposit, err := t.cashout(balance, transact)
		if err != nil {
			log.Println(err)
			if transactionStatus != 0 {
				resp = helpers.GetAPIResponse(transactionStatus, nil)
				return
			}
			resp = helpers.GetAPIResponse(consts.APIErrorUnknown, nil)
			return
		}
		err = t.BalanceRepository.Store(deposit)
		if err != nil {
			log.Println("Error when store balance", err)
			resp = helpers.GetAPIResponse(consts.APIErrorUnknown, nil)
		}
	case consts.Deposit:
		deposit, err := t.deposit(balance, transact)
		if err != nil {
			log.Println(err)
			resp = helpers.GetAPIResponse(consts.APIErrorUnknown, nil)
			return
		}
		err = t.BalanceRepository.Store(deposit)
		if err != nil {
			log.Println("Error when store balance", err)
			resp = helpers.GetAPIResponse(consts.APIErrorUnknown, nil)
			return
		}
	}
	resp = helpers.GetAPIResponse(consts.APIGeneralSuccess, nil)
	return
}

func (t *BalanceService) deposit(balance models.Balance, transact models.Transac) (deposit *models.Deposits, err error) {
	finalBalance := balance.Balance + transact.Amount
	var transactionHist models.TransactionHistory
	transactionHist.AccountID = transact.AccountID
	transactionHist.TransactionType = consts.Deposit
	transactionHist.FinalBalance = finalBalance
	transactionHist.TransactionAmount = transact.Amount

	lastId, err := t.TransactionHist.Store(transactionHist)
	if err != nil {
		return nil, err
	}

	deposit.AccountID = transact.AccountID
	deposit.Balance = finalBalance
	deposit.LastTransactionID = lastId

	return deposit, err
}

func (t *BalanceService) cashout(balance models.Balance, transact models.Transac) (transactionStatus int, deposit *models.Deposits, err error) {

	finalBalance := balance.Balance - transact.Amount
	if finalBalance < 0 {
		return consts.NotEnoughBalance, nil, fmt.Errorf("Balance not enough for this transaction")
	}
	var transactionHist models.TransactionHistory
	transactionHist.AccountID = transact.AccountID
	transactionHist.TransactionType = consts.CashOut
	transactionHist.FinalBalance = finalBalance
	transactionHist.TransactionAmount = transact.Amount

	lastId, err := t.TransactionHist.Store(transactionHist)
	if err != nil {
		return 0, nil, err
	}

	deposit.AccountID = transact.AccountID
	deposit.Balance = finalBalance
	deposit.LastTransactionID = lastId

	return consts.TransactionSuccess, deposit, nil
}
