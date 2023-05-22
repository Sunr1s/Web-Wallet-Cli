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

func (s *TxService) SaveTx(tx trx.Transaction) (int, error) {
	return s.repo.SaveTx(tx)
}

func (s *TxService) GetTx(userKey string) ([]trx.Transaction, error) {
	transactions, err := s.repo.GetTx(userKey)
	if err != nil {
		return nil, err
	}

	reversed := reverseTransactions(transactions)
	formatTransactionSendersAndReceivers(reversed)

	// return first 7 or fewer transactions
	if len(reversed) > 7 {
		return reversed[:7], nil
	}
	return reversed, nil
}

// reverseTransactions returns a reversed copy of the transactions slice.
func reverseTransactions(transactions []trx.Transaction) []trx.Transaction {
	reversed := make([]trx.Transaction, len(transactions))
	for i, tx := range transactions {
		reversed[len(transactions)-1-i] = tx
	}
	return reversed
}

// formatTransactionSendersAndReceivers updates the S_Sender and S_Reciver fields of each transaction.
func formatTransactionSendersAndReceivers(transactions []trx.Transaction) {
	for i, tx := range transactions {
		transactions[i].S_Sender = formatAddress(tx.Sender)
		transactions[i].S_Reciver = formatAddress(tx.Reciver)
	}
}

// formatAddress formats an address string. If the address is less than 20 characters long, it's returned as is.
func formatAddress(address string) string {
	if len(address) < 20 {
		return address
	}
	return address[10:14] + "-" + address[len(address)-16:len(address)-12]
}
