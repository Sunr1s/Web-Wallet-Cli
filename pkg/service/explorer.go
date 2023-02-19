package service

import (
	model "github.com/Sunr1s/webclient"
	bc "github.com/Sunr1s/webclient/pkg/blockchain"
	"github.com/Sunr1s/webclient/pkg/repository"
)

type ExplorerService struct {
	repo repository.Explorer
}

func NewExplorerService(repo repository.Explorer) *ExplorerService {
	return &ExplorerService{repo: repo}
}

func (e *ExplorerService) GetBlockChain(node string) []model.Block {

	blocks := e.repo.GetChain(node)

	var (
		block []model.Block
		tx    []model.BlockTransaction
		txs   []bc.Transaction
	)

	for i := range blocks {

		txs = blocks[i].Transactions

		block = append(block, model.Block{
			Height:      uint32(i),
			CurrHash:    bc.Base64Encode(blocks[i].CurrHash),
			CurrHashSrt: bc.Base64Encode(blocks[i].CurrHash)[:4] + "-" + bc.Base64Encode(blocks[i].CurrHash)[len(bc.Base64Encode(blocks[i].CurrHash))-4:],
			Miner:       blocks[i].Miner,
			MinerSrt:    blocks[i].Miner[10:14] + "-" + blocks[i].Miner[len(blocks[i].Miner)-16:len(blocks[i].Miner)-12],
			TimeStamp:   blocks[i].TimeStamp,
		})
		for j := range txs {
			block[i].Output += txs[j].Value
			tx = append(tx, model.BlockTransaction{
				CurrHash:    bc.Base64Encode(txs[j].CurrHash),
				CurrHashSrt: bc.Base64Encode(txs[j].CurrHash)[:4] + "-" + bc.Base64Encode(txs[j].CurrHash)[len(bc.Base64Encode(txs[j].CurrHash))-4:],
				Output:      txs[j].Value,
				TimeStamp:   blocks[i].TimeStamp,
			})
		}
		block[i].Transaction = tx

	}
	return block
}
