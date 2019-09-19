package repositories

import (
	"github.com/tnistest/internal/models"
	"github.com/tnistest/pkg/mysql"
)

type IAccountRepository interface {
	Store(account models.Accounts) (lastId int64, err error)
}

type AccountRepository struct {
	DB mysql.MySqlFactory
}

func NewAccountRepository(db mysql.MySqlFactory) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (a *AccountRepository) Store(account models.Accounts) (lastId int64, err error) {
	db, err := a.DB.GetDB()
	if err != nil {
		return
	}
	tx, err := db.Begin()

	if err != nil {
		return 0, err
	}
	if account.ID > 0 {

		stmtUptd, err := tx.Prepare(`update accounts set name =?, address=?, identity_number=?, identity_type=? where id=?`)
		_, err = stmtUptd.Exec(account.Name, account.Address, account.IdentityNumber, account.IdentityType, account.ID)

		if err != nil {
			if errTx := tx.Rollback(); errTx != nil {
				return 0, errTx
			}
			return 0, err
		}
		err = tx.Commit()
	} else {

		stmtInst, err := tx.Prepare(`insert into accounts(account_number, name, address, identity_number, identity_type) values (?,?,?,?,?)`)
		res, err := stmtInst.Exec(account.AccountNumber, account.Name, account.Address, account.IdentityNumber, account.IdentityType)

		if err != nil {
			if errTx := tx.Rollback(); errTx != nil {
				return 0, errTx
			}
			return 0, err
		}
		lastId, err = res.LastInsertId()
		err = tx.Commit()
		if err != nil {
			return 0, err
		}
	}
	return lastId, nil
}
