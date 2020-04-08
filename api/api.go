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
	a.r.HandleFunc("/", a.h.GetHome())
	a.r.HandleFunc("/register", a.h.RegisterUser()).Methods(http.MethodPost)
	a.r.HandleFunc("/login", a.h.Login()).Methods(http.MethodPost)
}
