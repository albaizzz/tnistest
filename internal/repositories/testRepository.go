package repositories

import (
	"github.com/tnistest/pkg/mysql"
)

type ITestRepository interface {
	Test() error
}

type TestRepository struct {
	db mysql.MySqlFactory
}

func NewTestRepository(mysqlDB mysql.MySqlFactory) *TestRepository {
	return &TestRepository{db: mysqlDB}
}

func (t *TestRepository) Test() error {
	return nil
}
