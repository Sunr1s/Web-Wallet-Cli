package handler

import (
	"fmt"
	"html/template"
	"net/http"

	model "github.com/Sunr1s/webclient"
	log "github.com/sirupsen/logrus"
)

type Explorer struct {
	Title    string
	Username string
	Chain    []model.Block
	Trx      []model.BlockTransaction
}

func (h *Handler) explorerPage(w http.ResponseWriter, r *http.Request) {

	var tx []model.BlockTransaction

	t, err := template.ParseFiles(TMPL_PATH+"base.gohtml",
		TMPL_PATH+"topbar.gohtml",
		TMPL_PATH+"footer.gohtml",
		TMPL_PATH+"navbar.gohtml",
		TMPL_PATH+"explorer.gohtml")
	if err != nil {
		fmt.Println("Error while parsing file", err)
		return
	}

	blocks := h.services.Explorer.GetBlockChain(":8088")

	for j := range blocks {
		for i := range blocks[j].Transaction {
			tx = append(tx, blocks[j].Transaction[i])
		}
	}

	_, Username, _ := getUser(r.Context())

	if err != nil {
		log.Error(err)
	}

	if len(tx) > 7 {
		tx = tx[:7]
	}

	p := Explorer{"Wellcome", Username, blocks, tx}

	//blocks[0].Transactions.
	err = t.Execute(w, p)

	if err != nil {
		fmt.Println("Error while executing template" + err.Error())
		return
	}
}
