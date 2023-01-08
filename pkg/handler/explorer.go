package handler

import (
	"fmt"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (h *Handler) explorerPage(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles(TMPL_PATH+"base.gohtml",
		TMPL_PATH+"topbar.gohtml",
		TMPL_PATH+"footer.gohtml",
		TMPL_PATH+"navbar.gohtml",
		TMPL_PATH+"explorer.gohtml")
	if err != nil {
		fmt.Println("Error while parsing file", err)
		return
	}

	Id, Username, Client := getUser(r.Context())

	log.Print(Client)

	if err != nil {
		log.Error(err)
	}

	p := Page{"Wellcome", Username, "", Id, 0, "", nil}

	err = t.Execute(w, p)

	if err != nil {
		fmt.Println("Error while executing template")
		return
	}
}
