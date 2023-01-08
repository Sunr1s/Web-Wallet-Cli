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

	// var (
	// 	res         *nt.Package
	// 	blocks      []*bc.Block
	// 	transacions [][]bc.Transaction
	// 	id          int
	// )

	// for i := 0; ; i++ {
	// 	res = nt.Send(":8088", &nt.Package{
	// 		Option: GET_BLOCK,
	// 		Data:   fmt.Sprintf("%d", i),
	// 	})
	// 	if res == nil || res.Data == "" {
	// 		break
	// 	}
	// 	blocks = append(blocks, bc.DeserializeBlock(res.Data))

	// }
	// for id := range blocks {
	// 	transacions = append(transacions, blocks[id].Transactions)
	// }
	// for id_block := range transacions {
	// 	for id_transac := range transacions[id_block] {
	// 		tx := transacions[id_block][id_transac]
	// 		if tx.Reciver == UserKey {
	// 			query := fmt.Sprintf("INSERT INTO %s (Reciver, Sender, Value, Message) values ($1, $2, $3, $4) RETURNING tx_id", txTable)

	// 			row := t.db.QueryRow(query, tx.Reciver, tx.Sender, tx.Value, tx.Signature)

	// 			if err := row.Scan(&id); err != nil {
	// 				return nil, err
	// 			}
	// 		} else {
	// 			continue
	// 		}
	// 	}
	// }

	rows, err := t.db.Query("SELECT tx_id, Reciver, Sender, Value, Message FROM transactions WHERE Sender =$1 OR Reciver =$1 ", UserKey)

	if err != nil {
		return nil, err
	}
	got := []trx.Transaction{}
	for rows.Next() {
		var r trx.Transaction
		err = rows.Scan(&r.Id, &r.Reciver, &r.Sender, &r.Value, &r.Message)
		if err != nil {
			return nil, err
		}
		got = append(got, r)
	}

	return got, nil
}
