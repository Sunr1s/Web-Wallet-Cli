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
