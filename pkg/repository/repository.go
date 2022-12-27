package repository

import (
	"database/sql"

	s_user "github.com/Sunr1s/webclient"
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

type Repository struct {
	Authorization
	Client
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSqlite(db),
		Client:        NewAuthSqlite(db),
	}
}
