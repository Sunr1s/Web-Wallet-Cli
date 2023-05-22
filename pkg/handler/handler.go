package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Sunr1s/webclient/pkg/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler struct {
	services *service.Service
}

var (
	Username string
	Id       int
)

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	workDir, err := os.Getwd()
	if err != nil {
		return nil
	}
	filesDir := filepath.Join(workDir, "static")
	h.HandleFileServer(r, "/files", http.Dir(filesDir))

	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", h.loginPage)
		r.Post("/login-submit", h.loginSubmit)

		r.Get("/register", h.signUpPage)
		r.Post("/reg-submit", h.signUpSubmit)
	})

	r.Route("/user", func(r chi.Router) {
		r.Use(h.userIdentity)

		r.Get("/wallet", h.walletPage)
		r.Post("/walletSubmit", h.walletSubmit)

		r.Get("/explorer", h.explorerPage)

		r.Get("/mywallet", h.myWalletPage)
		r.Get("/logout", h.logout)
	})

	r.Group(func(r chi.Router) {
		r.Use(h.pageIdentity)

		r.Get("/", h.mainPage)
	})

	return r
}

func (h *Handler) HandleFileServer(router chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		log.Errorf("FileServer does not permit any URL parameters.")
		return
	}

	if path != "/" && path[len(path)-1] != '/' {
		router.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	router.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
