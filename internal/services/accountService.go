package services

import (
	"log"

	"github.com/spf13/cast"

	"github.com/tnistest/internal/consts"
	"github.com/tnistest/internal/models"
	"github.com/tnistest/internal/repositories"
	"github.com/tnistest/pkg/helpers"
)

type IAccountService interface {
	Register(account models.Accounts) (resp models.APIResponse)
}

type AccountService struct {
	AccountRepo repositories.IAccountRepository
	BalanceRepo repositories.IBalanceRepository
}

func NewAccountService(accountRepo repositories.IAccountRepository, depositRepo repositories.IBalanceRepository) *AccountService {
	return &AccountService{AccountRepo: accountRepo, BalanceRepo: depositRepo}
}

func (a *AccountService) Register(account models.Accounts) (resp models.APIResponse) {
	lastId, err := a.AccountRepo.Store(account)
	if err != nil {
		log.Println(err)
		resp = helpers.GetAPIResponse(consts.APIErrorUnknown, nil)
		return
	}
	//store new deposit
	var deposit models.Deposits
	deposit.AccountID = cast.ToUint32(lastId)
	deposit.Balance = 0
	deposit.LastTransactionID = 0

	err = a.BalanceRepo.Store(&deposit)
	if err != nil {
		log.Println(err)
		resp = helpers.GetAPIResponse(consts.APIErrorUnknown, nil)
		return
	}
	resp = helpers.GetAPIResponse(consts.APIGeneralSuccess, nil)
	return
}
