package repository

import (
	"database/sql"
	"fmt"

	bc "github.com/Sunr1s/webclient/pkg/blockchain"
	nt "github.com/Sunr1s/webclient/pkg/network"
)

type ExplorerSqlite struct {
	db *sql.DB
}

func NewExplorerSqlite(db *sql.DB) *ExplorerSqlite {
	return &ExplorerSqlite{db: db}
}

func (e *ExplorerSqlite) GetChain(node string) []*bc.Block {
	var (
		res    *nt.Package
		blocks []*bc.Block
	)
	for i := 0; ; i++ {
		res = nt.Send(":8088", &nt.Package{
			Option: GET_BLOCK,
			Data:   fmt.Sprintf("%d", i),
		})
		if res == nil || res.Data == "" {
			break
		}
		blocks = append(blocks, bc.DeserializeBlock(res.Data))
	}
	return blocks
}
