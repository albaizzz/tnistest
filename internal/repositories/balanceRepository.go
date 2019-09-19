package repositories

import (
	"database/sql"

	"github.com/tnistest/internal/models"
	"github.com/tnistest/pkg/mysql"
)

type IBalanceRepository interface {
	GetByAccountNumber(accountNumber string) (models.Balance, error)
	GetById(accountID uint32) (models.Balance, error)
	Store(deposit *models.Deposits) error
}

type BalanceRepository struct {
	DB mysql.MySqlFactory
}

func NewBalanceRepository(db mysql.MySqlFactory) *BalanceRepository {
	return &BalanceRepository{DB: db}
}

func (b *BalanceRepository) GetByAccountNumber(accountNumber string) (balance models.Balance, err error) {
	db, err := b.DB.GetDB()
	if err != nil {
		return
	}

	rows, err := db.Query(`select d.balance, a.account_number from deposit d inner join accounts a on d.account_id = a.id
		where a.account_number = ?`, accountNumber)

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

func (b *BalanceRepository) GetById(accountID uint32) (balance models.Balance, err error) {
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

func (b *BalanceRepository) Store(deposit *models.Deposits) (err error) {
	db, err := b.DB.GetDB()
	if err != nil {
		return
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if deposit.LastTransactionID == 0 {

		stmtInst, err := tx.Prepare(`insert into deposits (account_id, balance, last_transaction_id) values (?,?,?)`)
		_, err = stmtInst.Exec(deposit.AccountID, deposit.Balance, deposit.LastTransactionID)

		if err != nil {
			if errTx := tx.Rollback(); errTx != nil {
				return errTx
			}
			return err
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
	} else {
		stmtInst, err := tx.Prepare(`update deposits set balance=?, last_transaction_id= ? where account_id =?`)
		_, err = stmtInst.Exec(deposit.Balance, deposit.LastTransactionID, deposit.AccountID)

		if err != nil {
			if errTx := tx.Rollback(); errTx != nil {
				return errTx
			}
			return err
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}
