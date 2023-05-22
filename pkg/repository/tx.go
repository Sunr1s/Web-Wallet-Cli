package repository

import (
	"database/sql"
	"fmt"

	trx "github.com/Sunr1s/webclient"
)

type TxSqlite struct {
	db *sql.DB
}

func NewTxSqlite(db *sql.DB) *TxSqlite {
	return &TxSqlite{db: db}
}

func (t *TxSqlite) SaveTx(tx trx.Transaction) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (Reciver, Sender, Value, Message) values ($1, $2, $3, $4) RETURNING tx_id", txTable)

	row := t.db.QueryRow(query, tx.Reciver, tx.Sender, tx.Value, tx.Message)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (t *TxSqlite) GetTx(UserKey string) ([]trx.Transaction, error) {
	rows, err := t.db.Query("SELECT tx_id, Reciver, Sender, Value, Message FROM transactions WHERE Sender =$1 OR Reciver =$1 ", UserKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []trx.Transaction
	for rows.Next() {
		var r trx.Transaction
		if err := rows.Scan(&r.Id, &r.Reciver, &r.Sender, &r.Value, &r.Message); err != nil {
			return nil, err
		}
		transactions = append(transactions, r)
	}

	// Check for errors from rows.Next() or rows.Scan().
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
