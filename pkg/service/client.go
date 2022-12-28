package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Sunr1s/webclient/pkg/blockchain"
	bc "github.com/Sunr1s/webclient/pkg/blockchain"
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

func (c *ClientService) ChainTX(spend, address string, User *blockchain.User) string {

	var ret string = ""
	repository.Address = []string{":8088", "9099"}
	num, err := strconv.Atoi(spend)

	if err != nil {
		log.Println("failed: strconv.Atoi(num)")
		return "failed: strconv.Atoi(num)"
	}

	for _, addr := range repository.Address {
		res := nt.Send(addr, &nt.Package{
			Option: repository.GET_LHASH,
		})
		if res == nil {
			continue
		}
		tx := bc.NewTransaction(User, bc.Base64Decode(res.Data), address, uint64(num))
		res = nt.Send(addr, &nt.Package{
			Option: repository.ADD_TRNSX,
			Data:   bc.SerializeTx(tx),
		})
		if res == nil {
			continue
		}
		if res.Data == "ok" {
			log.Printf("ok: (%s)\n", addr)
			ret += "ok: " + addr
		} else {
			log.Printf("fail: (%s) (%s)\n", addr, strings.Split(res.Data, "="))
			ret += fmt.Sprintf("fail: (%s) (%s)\n", addr, strings.Split(res.Data, "="))
		}
	}
	return ret
}
