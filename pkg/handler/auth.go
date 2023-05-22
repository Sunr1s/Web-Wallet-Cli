package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	s_user "github.com/Sunr1s/webclient"
	bc "github.com/Sunr1s/webclient/pkg/blockchain"
	log "github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId   int
	UserName string
	Client   *bc.User
}

type userData struct {
	Error    string
	Username string
	Id       int
}

func (h *Handler) loginPage(w http.ResponseWriter, r *http.Request) {
	data := &userData{}
	t := h.parseTemplates("base.gohtml", "topbar.gohtml", "footer.gohtml", "navbar.gohtml", "login.gohtml")
	if data.Error == "sql: no rows in result set" {
		data.Error = "Wrong username or password!"
	}
	data.Username = Username
	data.Id = Id
	h.executeTemplate(w, t, data)
}

func (h *Handler) loginSubmit(w http.ResponseWriter, r *http.Request) {
	data := &userData{}
	token, err := h.services.Authorization.GenerateToken(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		log.Error(err)
		data.Error = err.Error()
		http.Redirect(w, r, login, http.StatusSeeOther)
		return
	}

	log.Println(token)
	setCookieHandler(w, token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) signUpPage(w http.ResponseWriter, r *http.Request) {
	data := &userData{}
	t := h.parseTemplates("base.gohtml", "topbar.gohtml", "footer.gohtml", "navbar.gohtml", "signup.gohtml")
	if data.Error == "UNIQUE constraint failed: users.Username" {
		data.Error = "Username already taken"
	}
	h.executeTemplate(w, t, data)
}

func (h *Handler) signUpSubmit(w http.ResponseWriter, r *http.Request) {
	data := &userData{}
	var buf bytes.Buffer
	input := s_user.User{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}

	file, _, _ := r.FormFile("wallet")
	if r.FormValue("password") != r.FormValue("password_check") {
		data.Error = "Passwords don't match"
		log.Error(data.Error)
		http.Redirect(w, r, "/auth/register", http.StatusSeeOther)
		return
	}

	io.Copy(&buf, file)
	if buf.String() == "" {
		data.Error = "Error while parsing wallet"
		log.Error(data.Error)
		http.Redirect(w, r, "/auth/register", http.StatusSeeOther)
		return
	} else {
		input.Wallet = buf.String()
	}

	id, err := h.services.Authorization.RegisterUser(input)
	if err != nil {
		log.Error("Create user error")
		data.Error = err.Error()
		http.Redirect(w, r, "/auth/register", http.StatusSeeOther)
		return
	}

	log.Println("New user with id %d", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	clearCookieHandler(w)
	http.Redirect(w, r, login, http.StatusSeeOther)
}

func clearCookieHandler(w http.ResponseWriter) {
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

func (h *Handler) parseTemplates(filenames ...string) *template.Template {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("%s%s", TMPL_PATH, file))
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func (h *Handler) executeTemplate(w http.ResponseWriter, tmpl *template.Template, data *userData) {
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Error("Error while executing template" + err.Error())
	}
}
