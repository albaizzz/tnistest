package repositories

import (
	"database/sql"

	"github.com/tnistest/internal/models"
	"github.com/tnistest/pkg/mysql"
)

type ITransactionRepository interface {
	Store(transactionHistory models.TransactionHistory) (lastTransID int64, err error)
	Get(accountNumber string) []models.TransactionHistory
}

type TransactionRepository struct {
	DB mysql.MySqlFactory
}

func NewTransactionHistory(db mysql.MySqlFactory) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (t *TransactionRepository) Store(transactionHistory models.TransactionHistory) (lastId int64, err error) {
	db, err := t.DB.GetDB()
	if err != nil {
		return
	}

	tx, err := db.Begin()
	//insert new transaction history

	stmtInst, err := tx.Prepare(`insert into transaction_history (transaction_amount, final_balance, transaction_type, officer_id, account_id)
	VALUES(?,?,?,?,?)`)
	res, err := stmtInst.Exec(transactionHistory.TransactionAmount, transactionHistory.FinalBalance, transactionHistory.TransactionType, transactionHistory.OfficerID, transactionHistory.AccountID)
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
	return lastId, nil
}

func (t *TransactionRepository) GetByAccount(accountNumber string) (transactionList []models.UserTransactionHistory, err error) {
	db, err := b.DB.GetDB()
	if err != nil {
		return
	}

	rows, err := db.Query(`select d.balance, a.account_number from deposit d inner join accounts a on d.account_id = a.id
		where d.account_id = ?`, accountNumber)

	if err != nil {
		return
	}
	for rows.Next() {
		var transHist models.UserTransactionHistory
		err = rows.Scan(
			&transHist.AccountNumber,
			&transHist.TransactionAmount,
			&transHist.TransactionType,
		)
		transactionList = append(transactionList, transHist)
	}
	if err != nil && err != sql.ErrNoRows {
		return
	}
	return balance, nil
}
