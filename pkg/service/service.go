package service

import (
	model "github.com/Sunr1s/webclient"
	s_user "github.com/Sunr1s/webclient"
	"github.com/Sunr1s/webclient/pkg/blockchain"
	"github.com/Sunr1s/webclient/pkg/repository"
)

type Service struct {
	Authorization
	Client
	Transaction
	Explorer
}

func NewService(repos *repository.Repository, nodes []string) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Client:        NewClientService(repos.Client, nodes),
		Transaction:   NewTxService(repos.Transaction),
		Explorer:      NewExplorerService(repos.Explorer),
	}
}

type Authorization interface {
	CreateUser(user s_user.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, string, error)
}

type Client interface {
	ClientBalance(useraddr string) int
	LoadClient(Id int) *blockchain.User
	ChainTX(spend, address string, User *blockchain.User) string
}

type Transaction interface {
	SaveTx(tx s_user.Transaction) (int, error)
	GetTx(UserKey string) ([]s_user.Transaction, error)
}

type Explorer interface {
	GetBlockChain(node string) []model.Block
}
