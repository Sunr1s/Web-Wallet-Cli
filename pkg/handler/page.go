package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	TMPL_PATH  = "../templates/"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type Page struct {
	Title    string
	Username string
	Address  string
	Id       int
	Balance  int
}

func mainPage(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles(TMPL_PATH+"base.gohtml",
		TMPL_PATH+"topbar.gohtml",
		TMPL_PATH+"footer.gohtml",
		TMPL_PATH+"navbar.gohtml",
		TMPL_PATH+"index.gohtml")
	if err != nil {
		fmt.Println("Error while parsing file", err)
		return
	}

	Id, Username = getUser(r.Context())
	p := Page{"Wellcome", Username, "", Id, 0}

	err = t.Execute(w, p)
	if err != nil {
		fmt.Println("Error while executing template")
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
	log.Println(data.Client.Address())
	Balance := h.services.Client.ClientBalance(data.Client.Address())
	if Balance < 0 {
		log.Error("error conntect to any node's")
		Balance = 0
	}

	p := Page{"Wallet", Username, "", Id, Balance}

	err = t.Execute(w, p)

	if err != nil {
		fmt.Println("Error while executing template")
		return
	}
}

func (h *Handler) mywalletPage(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles(TMPL_PATH+"base.gohtml",
		TMPL_PATH+"topbar.gohtml",
		TMPL_PATH+"footer.gohtml",
		TMPL_PATH+"navbar.gohtml",
		TMPL_PATH+"mywallet.gohtml")
	if err != nil {
		fmt.Println("Error while parsing file", err)
		return
	}
	log.Println(data.Client.Address())
	Balance := h.services.Client.ClientBalance(data.Client.Address())
	Address := data.Client.Address()

	if Balance < 0 {
		log.Error("error conntect to any node's")
		Balance = 0
	}

	p := Page{"Wallet", Username, Address, Id, Balance}

	err = t.Execute(w, p)

	if err != nil {
		fmt.Println("Error while executing template")
		return
	}

}
