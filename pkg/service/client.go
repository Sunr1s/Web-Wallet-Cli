package service

import (
	"fmt"
	"strconv"

	"github.com/Sunr1s/webclient/pkg/blockchain"
	bc "github.com/Sunr1s/webclient/pkg/blockchain"
	nt "github.com/Sunr1s/webclient/pkg/network"
	"github.com/Sunr1s/webclient/pkg/repository"
)

type ClientService struct {
	repo  repository.Client
	nodes []string
}

func NewClientService(repo repository.Client, nodes []string) *ClientService {
	return &ClientService{repo: repo, nodes: nodes}
}

func (s *ClientService) LoadClient(clientID int) *blockchain.User {
	return s.repo.LoadClient(clientID)
}

func (s *ClientService) ClientBalance(userAddress string) (int, error) {
	for _, nodeAddress := range s.nodes {
		response := s.sendPackage(nodeAddress, repository.GET_BLNCE, userAddress)
		if response == nil {
			continue
		}
		balance, err := strconv.Atoi(response.Data)
		if err != nil || balance == 0 {
			continue
		}
		return balance, nil
	}
	return 0, fmt.Errorf("failed to retrieve balance for user %s", userAddress)
}

func (s *ClientService) ChainTX(amount, recipientAddress string, user *blockchain.User) (string, error) {
	value, err := strconv.Atoi(amount)
	if err != nil {
		return "", fmt.Errorf("failed to convert amount to integer: %w", err)
	}

	var result string
	for _, nodeAddress := range s.nodes {
		lastHashResponse := s.sendPackage(nodeAddress, repository.GET_LHASH, "")
		if lastHashResponse == nil {
			continue
		}
		lastHash, err := bc.Base64Decode(lastHashResponse.Data)
		if err != nil {
			continue
		}
		transaction, err := bc.NewTransaction(user, lastHash, recipientAddress, uint64(value))
		if err != nil {
			continue
		}
		transactionData, err := bc.SerializeTx(transaction)
		if err != nil {
			continue
		}
		response := s.sendPackage(nodeAddress, repository.ADD_TRNSX, transactionData)
		if response == nil || response.Data != "ok" {
			result += fmt.Sprintf("failure at node %s: %s\n", nodeAddress, response.Data)
			continue
		}
		result += fmt.Sprintf("success at node %s\n", nodeAddress)
	}
	if result == "" {
		return "", fmt.Errorf("transaction failed at all nodes")
	}
	return result, nil
}

func (s *ClientService) sendPackage(address string, option int, data string) *nt.Package {
	return nt.Send(address, &nt.Package{
		Option: option,
		Data:   data,
	})
}
