package service

import (
	trx "github.com/Sunr1s/webclient"
	"github.com/Sunr1s/webclient/pkg/repository"
)

type TxService struct {
	repo repository.Transaction
}

func NewTxService(repo repository.Transaction) *TxService {
	return &TxService{repo: repo}
}

func (t *TxService) SaveTx(tx trx.Transaction) (int, error) {
	return t.repo.SaveTx(tx)
}
