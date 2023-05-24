package service

import (
	"fmt"
	"time"

	model "github.com/Sunr1s/webclient"
	bc "github.com/Sunr1s/webclient/pkg/blockchain"
	"github.com/Sunr1s/webclient/pkg/repository"
)

type ExplorerService struct {
	repo repository.Explorer
}

// NewExplorerService initializes a new ExplorerService instance.
func NewExplorerService(repo repository.Explorer) *ExplorerService {
	return &ExplorerService{repo: repo}
}

// GetBlockChain retrieves the blockchain from a given node.
func (service *ExplorerService) GetBlockChain(node string) ([]model.Block, error) {
	blocks, err := service.repo.GetChain(node)
	if err != nil {
		return nil, err
	}

	modelBlocks := make([]model.Block, len(blocks))

	for i, block := range blocks {
		modelBlock, err := service.blockToModel(block, uint32(i))
		if err != nil {
			return nil, fmt.Errorf("failed to convert block to model: %w", err)
		}

		modelBlocks[len(blocks)-1-i] = modelBlock
	}

	return modelBlocks, nil
}

// blockToModel transforms a blockchain block to a model block.
func (service *ExplorerService) blockToModel(block *bc.Block, height uint32) (model.Block, error) {

	modelBlock := model.Block{
		PrevHash:    bc.Base64Encode(block.PrevHash),
		Height:      height,
		CurrHash:    bc.Base64Encode(block.CurrHash),
		Nonce:       block.Nonce,
		CurrHashSrt: formatHash(bc.Base64Encode(block.CurrHash)),
		Miner:       block.Miner,
		MinerSrt:    formatMiner(block.Miner),
		TimeStamp:   convertTimestampToString(block.TimeStamp),
		Output:      service.calculateTotalOutput(block.Transactions),
		Difficulty:  block.Difficulty,
		Signature:   bc.Base64Encode(block.Signature),
	}

	modelTransactions, err := service.transactionsToModel(block.Transactions, block)
	if err != nil {
		return model.Block{}, fmt.Errorf("failed to convert transactions to model: %w", err)
	}

	modelBlock.Transaction = modelTransactions

	return modelBlock, nil
}

// transactionsToModel transforms a slice of blockchain transactions to a slice of model transactions.
func (service *ExplorerService) transactionsToModel(transactions []bc.Transaction, block *bc.Block) ([]model.BlockTransaction, error) {
	modelTransactions := make([]model.BlockTransaction, len(transactions))

	for i, tx := range transactions {
		hash := bc.Base64Encode(tx.CurrHash)
		if hash == "" {
			hash = "CHAIN-REWARD"
		}

		modelTransactions[i] = model.BlockTransaction{
			CurrHash:    hash,
			CurrHashSrt: formatHash(hash),
			Output:      tx.Value,
			TimeStamp:   convertTimestampToString(block.TimeStamp),
			RandBytes:   bc.Base64Encode(tx.RandBytes),
			PrevBlock:   bc.Base64Encode(tx.PrevBlock),
			Sender:      tx.Sender,
			Receiver:    tx.Receiver,
			Value:       tx.Value,
			ToStorage:   tx.ToStorage,
			Signature:   bc.Base64Encode(tx.Signature),
		}
	}

	return modelTransactions, nil
}

// calculateTotalOutput calculates the total output of a slice of transactions.
func (service *ExplorerService) calculateTotalOutput(transactions []bc.Transaction) uint64 {
	var output uint64

	for _, tx := range transactions {
		output += tx.Value
	}

	return output
}

// formatHash returns a formatted string of a hash.
func formatHash(hash string) string {
	if len(hash) < 4 {
		return hash
	}
	return hash[:4] + "-" + hash[len(hash)-4:]
}

func convertTimestampToString(timestamp string) string {
	layout := "2006-01-02T15:04:05-07:00" // Layout for the given timestamp format
	t, err := time.Parse(layout, timestamp)
	if err != nil {
		return ""
	}

	// Format the time as a string in your desired format
	return t.Format("2006-01-02 15:04:05")
}

// formatMiner returns a formatted string of a miner.
func formatMiner(miner string) string {
	if len(miner) < 16 {
		return miner
	}
	return miner[10:14] + "-" + miner[len(miner)-16:len(miner)-12]
}
