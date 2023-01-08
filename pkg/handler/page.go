package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	trx "github.com/Sunr1s/webclient"
	log "github.com/sirupsen/logrus"
)

const (
	TMPL_PATH  = "../templates/"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type Page struct {
	Title        string
	Username     string
	Address      string
	Id           int
	Balance      int
	Error        string
	Transactions []trx.Transaction
}

func (h *Handler) mainPage(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		tx  []trx.Transaction
	)
	t, err := template.ParseFiles(TMPL_PATH+"base.gohtml",
		TMPL_PATH+"topbar.gohtml",
		TMPL_PATH+"footer.gohtml",
		TMPL_PATH+"navbar.gohtml",
		TMPL_PATH+"index.gohtml")
	if err != nil {
		fmt.Println("Error while parsing file", err)
		return
	}

	Id, Username, Client := getUser(r.Context())

	if Client != nil {
		tx, err = h.services.GetTx(Client.Address())
	}

	if err != nil {
		log.Error(err)
	}

	for i := range tx {
		tx[i].S_Sender = tx[i].Sender[10:14] + "-" + tx[i].Sender[len(tx[i].Sender)-16:len(tx[i].Sender)-12]
		tx[i].S_Reciver = tx[i].Reciver[10:14] + "-" + tx[i].Reciver[len(tx[i].Reciver)-16:len(tx[i].Sender)-12]
	}

	p := Page{"Wellcome", Username, "", Id, 0, "", tx}

	err = t.Execute(w, p)
	if err != nil {
		fmt.Println("Error while executing template" + err.Error())
		return
	}
}

func (h *Handler) walletPage(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles(TMPL_PATH+"base.gohtml",
		TMPL_PATH+"topbar.gohtml",
		TMPL_PATH+"footer.gohtml",
		TMPL_PATH+"navbar.gohtml",
		TMPL_PATH+"wallet.gohtml")
	if err != nil {
		fmt.Println("Error while parsing file", err)
		return
	}

	Id, Username, Client := getUser(r.Context())
	Balance := h.services.Client.ClientBalance(Client.Address())

	if Balance < 0 {
		log.Error("error conntect to any node's")
		Balance = 0
	}

	p := Page{"Wallet", Username, "", Id, Balance, data.Error, nil}
	err = t.Execute(w, p)

	if err != nil {
		fmt.Println("Error while executing template")
		return
	}
}

func (h *Handler) walletSubmit(w http.ResponseWriter, r *http.Request) {

	_, _, Client := getUser(r.Context())

	input := trx.Transaction{
		Reciver: r.FormValue("reciver"),
		Value:   r.FormValue("amount"),
		Message: r.FormValue("message"),
		Sender:  Client.Address(),
	}

	out := h.services.Client.ChainTX(input.Value, input.Reciver, Client)

	log.Println(out)

	if out == "" || strings.Contains(out, "fail") {
		data.Error = out
		http.Redirect(w, r, "/user/wallet", http.StatusSeeOther)
	} else {
		id, err := h.services.Transaction.SaveTx(input)
		if id == 0 || err != nil {
			data.Error = err.Error()
			http.Redirect(w, r, "/user/wallet", http.StatusSeeOther)
		}
		http.Redirect(w, r, "/user/wallet", http.StatusSeeOther)
	}
}

func (h *Handler) myWalletPage(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles(TMPL_PATH+"base.gohtml",
		TMPL_PATH+"topbar.gohtml",
		TMPL_PATH+"footer.gohtml",
		TMPL_PATH+"navbar.gohtml",
		TMPL_PATH+"mywallet.gohtml")
	if err != nil {
		fmt.Println("Error while parsing file", err)
		return
	}

	Id, Username, Client := getUser(r.Context())

	Balance := h.services.Client.ClientBalance(Client.Address())
	Address := Client.Address()

	if Balance < 0 {
		log.Error("error conntect to any node's")
		Balance = 0
	}

	p := Page{"Wallet", Username, Address, Id, Balance, "", nil}

	err = t.Execute(w, p)

	if err != nil {
		fmt.Println("Error while executing template")
		return
	}

}
