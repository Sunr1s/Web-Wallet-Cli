package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sunr1s/webclient/pkg/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

var (
	Username string
	Id       int
)

func (h *Handler) InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	h.HandleFileServer(r, "/files", filesDir)

	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", h.loginPage)
		r.Post("/login-submit", h.loginSubmit)

		r.Get("/register", h.singUpPage)
		r.Post("/reg-submit", h.singUpSubmit)
	})
	r.Route("/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.userIdentity)
			r.Get("/wallet", h.walletPage)
			r.Get("/mywallet", h.mywalletPage)
			r.Get("/logout", h.logout)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(h.pageIdentity)
		r.Get("/", mainPage)
	})

	return r
}

func (h *Handler) HandleFileServer(router chi.Router, path string, root http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.ContainsAny(path, "{}*") {
			panic("FileServer does not permit any URL parameters.")
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
	})
}
