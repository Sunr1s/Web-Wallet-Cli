package handler

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	login               = "/auth/login"
)

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := getCookieHandler(w, r)

		if errors.Is(err, http.ErrNoCookie) {
			http.Redirect(w, r, login, http.StatusSeeOther)
		}

		header := r.Header.Get("Authorization")

		if header == "" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("empty auth header")
			return
		}
		log.Println(header)

		UserId, UserName, err := h.services.Authorization.ParseToken(header)

		data.Client = h.services.Client.LoadClient(UserId)

		newCtx := context.WithValue(r.Context(), userCtx, &tokenClaims{
			UserId:   UserId,
			UserName: UserName,
		})

		r = r.WithContext(newCtx)

		next.ServeHTTP(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err.Error())
			return
		}
	})
}

func (h *Handler) pageIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := getCookieHandler(w, r)
		var newCtx context.Context
		if errors.Is(err, http.ErrNoCookie) {
			newCtx = context.WithValue(r.Context(), userCtx, &tokenClaims{
				UserId:   0,
				UserName: "",
			})
		} else {
			header := r.Header.Get("Authorization")

			if header == "" {
				w.WriteHeader(http.StatusUnauthorized)
				log.Println("empty auth header")
				return
			}

			UserId, UserName, err := h.services.Authorization.ParseToken(header)
			newCtx = context.WithValue(r.Context(), userCtx, &tokenClaims{
				UserId:   UserId,
				UserName: UserName,
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				log.Println(err.Error())
				return
			}
		}
		r = r.WithContext(newCtx)
		next.ServeHTTP(w, r)
	})
}
func getUser(ctx context.Context) (int, string) {
	raw := ctx.Value(userCtx).(*tokenClaims)
	return raw.UserId, raw.UserName
}
