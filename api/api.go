package api

import (
	"github.com/gorilla/mux"
	"go-blog-api/api/handlers"
	"go-blog-api/config"
	"net/http"
)

type API struct {
	h *handlers.Handler
	r *Router
}

type Router struct {
	*mux.Router
}

func NewAPI(config *config.AppConfig) *API {
	return &API{
		h: &handlers.Handler{
			AppConfig: config,
		},
		r: &Router{mux.NewRouter()},
	}
}

func (a *API) initialize() {
	//unauthenticated routes
	publicRoutes := a.r

	publicRoutes.HandleFunc("/", a.h.GetHome())
	publicRoutes.HandleFunc("/register", a.h.RegisterUser()).Methods(http.MethodPost)
	publicRoutes.HandleFunc("/login", a.h.Login()).Methods(http.MethodPost)

	//Authenticated routes
	a.authRoutes("/articles", a.h.CreateArticle(), http.MethodPost)
	a.authRoutes("/article/{id}", a.h.GetOneArticle(), http.MethodGet)
	a.authRoutes("/articles", a.h.GetArticles(), http.MethodGet)
	a.authRoutes("/user/articles", a.h.GetLoggedInUserArticles(), http.MethodGet)
	a.authRoutes("/article/{id}", a.h.UpdateArticle(), http.MethodPut)
	a.authRoutes("/article/{id}", a.h.DeleteArticle(), http.MethodDelete)

	a.authRoutes("/article/{id}/comment", a.h.PostComment(), http.MethodPost)
	a.authRoutes("/comment/{id}", a.h.GetOneComment(), http.MethodGet)
	a.authRoutes("/article/{id}/comment", a.h.GetArticleComments(), http.MethodGet)
	a.authRoutes("/user/comments", a.h.GetUserComments(), http.MethodGet)
	a.authRoutes("/comment/{id}", a.h.UpdateComment(), http.MethodPut)
	a.authRoutes("/comment/{id}", a.h.DeleteComment(), http.MethodDelete)

	a.authRoutes("/article/{id}/like", a.h.AddOrDislikeAnArticle(), http.MethodPost)
	a.authRoutes("/comment/{id}/like", a.h.AddorDislikeAComment(), http.MethodPost)
}

func (a *API) authRoutes(path string, handler http.HandlerFunc, method string) {
	a.r.HandleFunc(path, AuthMiddleware(handler)).Methods(method)
}
