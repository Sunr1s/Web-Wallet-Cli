package service

import (
	s_user "github.com/Sunr1s/webclient"
	"github.com/Sunr1s/webclient/pkg/blockchain"
	"github.com/Sunr1s/webclient/pkg/repository"
)

type Service struct {
	Authorization
	Client
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Client:        NewClientService(repos.Client),
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
}
