package repository

import (
	"database/sql"

	s_user "github.com/Sunr1s/webclient"
	trx "github.com/Sunr1s/webclient"
	bc "github.com/Sunr1s/webclient/pkg/blockchain"
)

type Authorization interface {
	CreateUser(user s_user.User) (int, error)
	GetUser(username, password string) (s_user.User, error)
}

type Client interface {
	LoadClient(Id int) *bc.User
	ReadFile(filename string) string
}

type Transaction interface {
	SaveTx(tx s_user.Transaction) (int, error)
	GetTx(UserKey string) ([]trx.Transaction, error)
}

type Repository struct {
	Authorization
	Transaction
	Client
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSqlite(db),
		Client:        NewAuthSqlite(db),
		Transaction:   NewTxSqlite(db),
	}
}
