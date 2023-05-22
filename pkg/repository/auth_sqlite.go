package repository

import (
	"database/sql"
	"fmt"

	s_user "github.com/Sunr1s/webclient"
	bc "github.com/Sunr1s/webclient/pkg/blockchain"
)

// SQLiteAuthRepository encapsulates operations with SQLite database
type SQLiteAuthRepository struct {
	db *sql.DB
}

// NewSQLiteAuthRepository initializes a new SQLiteAuthRepository
func NewSQLiteAuthRepository(db *sql.DB) *SQLiteAuthRepository {
	return &SQLiteAuthRepository{db: db}
}

// RegisterUser inserts a new user into the database
func (r *SQLiteAuthRepository) RegisterUser(user s_user.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (Username, Password, Email, Wallet) values ($1, $2, $3, $4) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Username, user.Password, user.Email, user.Wallet)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// RetrieveUser fetches a user from the database based on username and password
func (r *SQLiteAuthRepository) GetUser(username, password string) (s_user.User, error) {
	query := "SELECT * FROM users WHERE Username = $1 AND Password = $2"
	row := r.db.QueryRow(query, username, password)

	var user s_user.User
	if err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Wallet); err != nil {
		return s_user.User{}, err
	}

	return user, nil
}

// LoadClient retrieves a user's wallet based on the user's id
func (r *SQLiteAuthRepository) LoadClient(id int) *bc.User {
	query := "SELECT Wallet FROM users WHERE Id = $1"
	row := r.db.QueryRow(query, id)

	var data string
	if err := row.Scan(&data); err != nil {
		return nil
	}

	privKey := readKeys(data, true)
	if privKey == "" {
		return nil
	}

	user := bc.LoadUser(privKey)
	if user == nil {
		return nil
	}

	return user
}
