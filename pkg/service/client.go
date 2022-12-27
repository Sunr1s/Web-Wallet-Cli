package service

import (
	"strconv"

	"github.com/Sunr1s/webclient/pkg/blockchain"
	nt "github.com/Sunr1s/webclient/pkg/network"
	"github.com/Sunr1s/webclient/pkg/repository"
)

type ClientService struct {
	repo repository.Client
}

func NewClientService(repo repository.Client) *ClientService {
	return &ClientService{repo: repo}
}

func (c *ClientService) LoadClient(Id int) *blockchain.User {
	Client := c.repo.LoadClient(Id)
	return Client
}

func (c *ClientService) ClientBalance(useraddr string) int {
	var (
		balance int
		err     error
	)
	repository.Address = []string{":8088", "9099"}

	for _, addr := range repository.Address {
		res := nt.Send(addr, &nt.Package{
			Option: repository.GET_BLNCE,
			Data:   useraddr,
		})
		if res == nil {
			return -1
		}
		balance, err = strconv.Atoi(res.Data)
		if err != nil && balance == 0 {
			return -1
		} else {
			return balance
		}
	}
	return balance
}
