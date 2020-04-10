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
	couldNotGetArticles   = "could not get all articles"
	articleEmpty = "article record is empty"
	actionDisallowed = "action is disallowed"
	couldnotupdateArticle = "could not update article"
	deleteArticleFailed = "delete article failed"
)

func (h *Handler) CreateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articleRequest := dal.Article{}

		err := unMarshalRequest(r, &articleRequest)
		if err != nil {
			http.Error(w, unableToUnMarshalRequest, http.StatusInternalServerError)
			return
		}

		userID, ok := isUserIDValid(w, r)
		if !ok {
			return
		}

		valErrors := isCreateArticleRequestValid(articleRequest)
		if len(valErrors) > 0 {
			writeValErrors(w, valErrors)
			return
		}

		articleRequest.ID = services.GenerateUUID()
		articleRequest.CreatedAt = time.Now()
		articleRequest.AuthorID = userID

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
		articleID, ok := isArticleIDValid(r, w)
		if !ok {
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

func isArticleIDValid(r *http.Request, w http.ResponseWriter) (string, bool) {
	articleID := mux.Vars(r)["id"]

	if !isValidGUID(articleID) {
		http.Error(w, invalidArticleID, http.StatusBadRequest)
		return "", false
	}
	return articleID, true
}

func (h *Handler) GetArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dao, failed := GetDAOFromDB(h, w)
		if failed {
			return
		}
		defer dao.DB.Close()

		articles, err := dao.GetAllArticles()
		if err != nil {
			http.Error(w, couldNotGetArticles, http.StatusInternalServerError)
			return
		}

		if len(articles) == 0 {
			http.Error(w, articleEmpty, http.StatusNotFound)
			return
		}

		writeResponse(w, articles)
	}
}

func( h *Handler) GetLoggedInUserArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dao, failed := GetDAOFromDB(h, w)
		if failed {
			return
		}
		defer dao.DB.Close()

		userID, ok := isUserIDValid(w, r)
		if !ok {
			return
		}

		articles, err := dao.GetArticlesByUserID(userID)
		if err != nil {
			http.Error(w, couldNotGetArticles, http.StatusInternalServerError)
			return
		}

		if len(articles) == 0 {
			http.Error(w, articleEmpty, http.StatusNotFound)
			return
		}

		writeResponse(w, articles)
	}
}

func (h *Handler) UpdateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articleRequest := dal.Article{}

		err := unMarshalRequest(r, &articleRequest)
		if err != nil {
			http.Error(w, unableToUnMarshalRequest, http.StatusInternalServerError)
			return
		}

		userID, ok := isUserIDValid(w, r)
		if !ok {
			return
		}

		article, dao, err, done := initializeArticleParamHandler(w, r, h)
		if done {
			return
		}
		defer dao.DB.Close()

		if article.AuthorID != userID {
			http.Error(w, actionDisallowed, http.StatusForbidden)
			return
		}

		articleRequest.ID = article.ID
		updArticle, err := dao.UpdateArticle(articleRequest)

		if err != nil {
			http.Error(w, couldnotupdateArticle, http.StatusInternalServerError)
			return
		}

		updArticle.AuthorID = article.AuthorID
		writeResponse(w, updArticle)
	}
}

func (h *Handler) DeleteArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := isUserIDValid(w, r)
		if !ok {
			return
		}

		article, dao, err, done := initializeArticleParamHandler(w, r, h)
		if done {
			return
		}
		defer dao.DB.Close()

		if article.AuthorID != userID {
			http.Error(w, actionDisallowed, http.StatusForbidden)
			return
		}

		err = dao.DeleteArticle(article.ID)
		if err != nil {
			http.Error(w, deleteArticleFailed, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func initializeArticleParamHandler(w http.ResponseWriter, r *http.Request, h *Handler) (*dal.Article, *dal.DAL, error, bool) {

	articleID, ok := isArticleIDValid(r, w)
	if !ok {
		return nil, nil, nil, true
	}

	dao, failed := GetDAOFromDB(h, w)
	if failed {
		return nil, nil, nil, true
	}

	article, err := dao.GetArticleById(articleID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			http.Error(w, articleEmpty, http.StatusNotFound)
			return nil, nil, nil, true
		}

		http.Error(w, articlelookupfailed, http.StatusInternalServerError)
		return nil, nil, nil, true
	}


	return article, dao, nil, false
}

func isUserIDValid(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID:= r.Header.Get("userID")
	if userID == "" {
		http.Error(w, userIDNil, http.StatusBadRequest)
		return "", false
	}

	return userID, true
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
