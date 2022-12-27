package repository

import (
	"database/sql"
	"fmt"

	s_user "github.com/Sunr1s/webclient"
	bc "github.com/Sunr1s/webclient/pkg/blockchain"
)

type AuthSqlite struct {
	db *sql.DB
}

func NewAuthSqlite(db *sql.DB) *AuthSqlite {
	return &AuthSqlite{db: db}
}

func (r *AuthSqlite) CreateUser(user s_user.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (Username, Password, Email, Wallet) values ($1, $2, $3, $4) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Username, user.Password, user.Email, user.Wallet)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthSqlite) GetUser(username, password string) (s_user.User, error) {
	var user s_user.User
	row := r.db.QueryRow("SELECT * FROM users WHERE Username= $1 AND Password = $2;", username, password)
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Wallet)
	return user, err
}

func (r *AuthSqlite) LoadClient(Id int) *bc.User {
	var data string

	row := r.db.QueryRow("SELECT Wallet FROM users WHERE Id = $1;", Id)
	err := row.Scan(&data)

	if err != nil {
		return nil
	}
	priv := readKeys(data, true)
	if priv == "" {
		return nil
	}
	user := bc.LoadUser(priv)
	if user == nil {
		return nil
	}
	return user
}
