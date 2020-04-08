package handlers

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"go-blog-api/dal"
	"go-blog-api/services"
	"net/http"
	"time"
)

const (
	userIDNil             = "user ID is not set"
	unabletoCreateArticle = "unable to create article"
	articlelookupfailed   = "article look up failed"
	noArticlefound        = "no article is found"
	invalidArticleID      = "invalid article id"
)

func (h *Handler) CreateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articleRequest := dal.Article{}

		err := unMarshalRequest(r, &articleRequest)
		if err != nil {
			http.Error(w, unableToUnMarshalRequest, http.StatusInternalServerError)
			return
		}

		if r.Header.Get("userID") == "" {
			http.Error(w, userIDNil, http.StatusInternalServerError)
			return
		}

		valErrors := isCreateArticleRequestValid(articleRequest)
		if len(valErrors) > 0 {
			writeValErrors(w, valErrors)
			return
		}

		articleRequest.ID = services.GenerateUUID()
		articleRequest.CreatedAt = time.Now()
		articleRequest.AuthorID = r.Header.Get("userID")

		dao, failed := GetDAOFromDB(h, w)
		if failed {
			return
		}
		defer dao.DB.Close()

		article, err := dao.CreateNewArticle(articleRequest)
		if err != nil {
			http.Error(w, unabletoCreateArticle, http.StatusInternalServerError)
			return
		}
		article.User.Password = ""
		w.WriteHeader(http.StatusCreated)
		writeResponse(w, article)
	}
}
func (h *Handler) GetOneArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articleID := mux.Vars(r)["id"]

		if !isValidGUID(articleID) {
			http.Error(w, invalidArticleID, http.StatusBadRequest)
			return
		}

		dao, failed := GetDAOFromDB(h, w)
		if failed {
			return
		}
		defer dao.DB.Close()

		article, err := dao.GetArticleById(articleID)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				http.Error(w, noArticlefound, http.StatusNotFound)
				return
			}

			http.Error(w, articlelookupfailed, http.StatusInternalServerError)
			return
		}

		article.User.Password = ""
		writeResponse(w, article)
	}
}

func GetDAOFromDB(h *Handler, w http.ResponseWriter) (*dal.DAL, bool) {
	db, err := services.NewDB(h.AppConfig)
	if err != nil {
		http.Error(w, unableToConnectToDB, http.StatusInternalServerError)
		return nil, true
	}

	dao := dal.NewDAL(db)
	return dao, false
}
