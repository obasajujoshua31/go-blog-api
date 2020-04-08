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

}

func (a *API) authRoutes(path string, handler http.HandlerFunc, method string) {
	a.r.HandleFunc(path, AuthMiddleware(handler)).Methods(method)
}
