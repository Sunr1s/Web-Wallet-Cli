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

func (t *TxService) GetTx(UserKey string) ([]trx.Transaction, error) {
	txs, err := t.repo.GetTx(UserKey)

	var reverse []trx.Transaction

	for i := len(txs) - 1; i >= 0; i-- {
		reverse = append(reverse, txs[i])
	}
	if reverse != nil {
		if len(reverse) > 7 {
			return reverse[:7], err
		} else {
			return reverse, err
		}
	} else {
		return nil, err
	}

}
