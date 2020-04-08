package api

import (
	"go-blog-api/services"
	"net/http"
)

const (
	token        = "token"
	noTokenFound = "no token found"
)

func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get(token)
		if token == "" {
			http.Error(w, noTokenFound, http.StatusUnauthorized)
			return
		}

		userId, err := services.GetUserIDFromToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		r.Header.Set("userID", userId)
		next.ServeHTTP(w, r)
	}
}
