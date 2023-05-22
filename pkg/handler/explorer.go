package handler

import (
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
	t := h.parseTemplates("base.gohtml", "topbar.gohtml", "footer.gohtml", "navbar.gohtml", "explorer.gohtml")

	blocks, err := h.services.Explorer.GetBlockChain(":8088")
	if err != nil {
		log.Error("Error getting blockchain: ", err)
		return
	}

	var tx []model.BlockTransaction
	for _, block := range blocks {
		tx = append(tx, block.Transaction...)
	}

	_, Username, _ := getUser(r.Context())

	if len(tx) > 7 {
		tx = tx[:7]
	}

	p := Explorer{"Wellcome", Username, blocks, tx}

	err = t.Execute(w, p)
	if err != nil {
		log.Error("Error executing template: ", err)
	}
}
