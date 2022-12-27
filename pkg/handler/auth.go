package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	s_user "github.com/Sunr1s/webclient"
	"github.com/Sunr1s/webclient/pkg/blockchain"

	log "github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId   int
	UserName string
}

var data struct {
	Error    string
	Username string
	Id       int
	Client   *blockchain.User
}

func (h *Handler) loginPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(TMPL_PATH+"base.gohtml",
		TMPL_PATH+"topbar.gohtml",
		TMPL_PATH+"footer.gohtml",
		TMPL_PATH+"navbar.gohtml",
		TMPL_PATH+"login.gohtml")
	if err != nil {
		fmt.Println("Error while parsing file", err)
		return
	}
	if data.Error == "sql: no rows in result set" {
		data.Error = "Wrong username or password!"
	}
	data.Username = Username
	data.Id = Id
	err = t.Execute(w, data)

	if err != nil {
		fmt.Println("Error while executing template" + err.Error())
		return
	}
}

func (h *Handler) loginSubmit(w http.ResponseWriter, r *http.Request) {

	//privkey := r.FormValue("privkey")
	token, err := h.services.Authorization.GenerateToken(
		r.FormValue("username"),
		r.FormValue("password"))

	if err != nil {
		log.Error(err)
		data.Error = err.Error()
		http.Redirect(w, r, login, 303)
	} else {
		log.Println(token)

		setCookieHandler(w, token)

		data.Error = ""
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}

}

func (h *Handler) singUpPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(TMPL_PATH+"base.gohtml",
		TMPL_PATH+"topbar.gohtml",
		TMPL_PATH+"footer.gohtml",
		TMPL_PATH+"navbar.gohtml",
		TMPL_PATH+"singup.gohtml")
	if err != nil {
		fmt.Println("Error while parsing file", err)
		return
	}

	if data.Error == "UNIQUE constraint failed: users.Username" {
		data.Error = "Username already taken"
	}

	err = t.Execute(w, data)

	if err != nil {
		fmt.Println("Error while executing template" + err.Error())
		return
	}
}

func (h *Handler) singUpSubmit(w http.ResponseWriter, r *http.Request) {

	var input s_user.User
	var buf bytes.Buffer

	input.Username = r.FormValue("username")
	input.Password = r.FormValue("password")
	input.Email = r.FormValue("email")

	file, _, _ := r.FormFile("wallet")

	if r.FormValue("password") != r.FormValue("password_check") {
		data.Error = "Passwords doesn't match"
		log.Error(data.Error)
		http.Redirect(w, r, "/auth/register", 303)
	}

	io.Copy(&buf, file)
	if buf.String() == "" {
		data.Error = "Error while parsing wallet"
		log.Error(data.Error)
		http.Redirect(w, r, "/auth/register", 303)
	} else {
		input.Wallet = buf.String()
	}

	buf.Reset()

	id, err := h.services.Authorization.CreateUser(input)

	if err != nil {
		data.Error = err.Error()
		http.Redirect(w, r, "/auth/register", 303)
	}

	if err != nil {
		log.Error("Create user error")
		data.Error = err.Error()
		http.Redirect(w, r, "/auth/register", 303)
	}

	log.Println("New user with id %d", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	r.Header.Add("Authorization", "")
	Username, Id = "", 0

	http.Redirect(w, r, login, http.StatusSeeOther)
}

func setCookieHandler(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", cookie.Value)
	return nil
}
