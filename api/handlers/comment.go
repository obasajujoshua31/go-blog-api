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
	couldNotcreateComment = "could not create comment"
	invalidCommentID = "invalid Comment id"
	couldNotGetComment = "could not get comment"
	couldNotGetComments = "could not get comments"
	noCommentFound = "no comment found"
	couldNotUpdateComment = "could not update comment"
	couldNotDeleteComment = "could not delete comment"
)

func (h *Handler) PostComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := isUserIDValid(w, r)
		if !ok {
			return
		}

		commentRequest, good := isRequestMarshalGood(w,  r)
		if !good {
			return
		}

		article, dao, err, done := initializeArticleParamHandler(w, r, h)
		if done {
			return
		}
		defer dao.DB.Close()

		valErrors := isCreateCommentValid(*commentRequest)
		if len(valErrors) > 0 {
			writeValErrors(w, valErrors)
			return
		}

		commentRequest.ID = services.GenerateUUID()
		commentRequest.CreatedAt = time.Now()
		commentRequest.ArticleID = article.ID
		commentRequest.ReviewerID = userID


		comment, err := dao.CreateComment(*commentRequest, *article)
		if err != nil {
			http.Error(w, couldNotcreateComment, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		comment.User.Password = ""
		writeResponse(w, comment)
	}
}

func (h *Handler) GetOneComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		commentID, ok := isCommentIDValid(w, r)

		if !ok{
			return
		}

		dao, failed := GetDAOFromDB(h, w)
		if failed {
			return
		}

		comment, err := dao.GetCommentByID(commentID)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				http.Error(w, noCommentFound, http.StatusNotFound)
				return
			}

			http.Error(w, couldNotGetComment, http.StatusBadRequest)
			return
		}

		comment.User.Password = ""
		writeResponse(w, comment)
	}
}

func (h *Handler) GetArticleComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		article, dao, _, done := initializeArticleParamHandler(w, r, h)
		if done {
			return
		}
		defer dao.DB.Close()

		comments, err := dao.GetCommentOnArticle(article.ID)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				http.Error(w, noCommentFound, http.StatusNotFound)
				return
			}

			http.Error(w, couldNotGetComments, http.StatusInternalServerError)
			return
		}

		writeResponse(w, comments)
	}
}

func (h *Handler) GetUserComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := isUserIDValid(w, r)
		if !ok {
			return
		}

		dao, failed := GetDAOFromDB(h, w)
		if failed {
			return
		}

		comments, err := dao.GetCommentByUserID(userID)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				http.Error(w, noCommentFound, http.StatusNotFound)
				return
			}

			http.Error(w, couldNotGetComments, http.StatusInternalServerError)
			return
		}

		writeResponse(w, comments)
	}
}

func (h *Handler) UpdateComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := isUserIDValid(w, r)
		if !ok {
			return
		}


		commentRequest, good := isRequestMarshalGood(w,  r)
		if !good {
			return
		}

		comment, dao, failed := initializeCommentParamHandler(w, r, h)
		if failed {
			return
		}

		if comment.ReviewerID != userID  {
			http.Error(w, actionDisallowed, http.StatusForbidden)
			return
		}

		valErrors := isCreateCommentValid(*commentRequest)
		if len(valErrors) > 0 {
			writeValErrors(w, valErrors)
			return
		}

		comment.Content = commentRequest.Content
		updComm, err := dao.UpdateComment(*comment)
		if err != nil {
			http.Error(w, couldNotUpdateComment, http.StatusInternalServerError)
			return
		}

		updComm.User.Password = ""
		writeResponse(w, updComm)
	}
}

func (h *Handler) DeleteComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := isUserIDValid(w, r)
		if !ok {
			return
		}

		comment, dao, failed := initializeCommentParamHandler(w, r, h)
		if failed {
			return
		}

		article, err := dao.GetArticleById(comment.ArticleID)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				http.Error(w, articleEmpty, http.StatusNotFound)
				return
			}


			http.Error(w, articlelookupfailed, http.StatusInternalServerError)
			return
		}

		if comment.ReviewerID != userID  {
			http.Error(w, actionDisallowed, http.StatusForbidden)
			return
		}

		err = dao.DeleteComment(comment.ID, *article)
		if err != nil {
			http.Error(w, couldNotDeleteComment, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func isCommentIDValid(w http.ResponseWriter, r *http.Request) (string, bool) {
	commentID := mux.Vars(r)["id"]

	if !isValidGUID(commentID) {
		http.Error(w, invalidCommentID, http.StatusBadRequest)
		return "", false
	}
	return commentID, true
}

func initializeCommentParamHandler(w http.ResponseWriter, r *http.Request, h *Handler) (*dal.Comment, *dal.DAL, bool) {

	commentID, ok := isCommentIDValid(w, r)
	if !ok {
		return nil, nil, true
	}

	dao, failed := GetDAOFromDB(h, w)
	if failed {
		return nil, nil, true
	}

	comment, err := dao.GetCommentByID(commentID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			http.Error(w, noCommentFound, http.StatusNotFound)
			return nil, nil, true
		}

		http.Error(w, couldNotGetComment, http.StatusInternalServerError)
		return nil, nil, true
	}

	return comment, dao, false
}

func isRequestMarshalGood(w http.ResponseWriter, r *http.Request) (*dal.Comment, bool) {
	commentRequest := dal.Comment{}
	err := unMarshalRequest(r, &commentRequest)
	if err != nil {
		http.Error(w, unableToUnMarshalRequest, http.StatusInternalServerError)
		return nil, false
	}

	return &commentRequest, true
}
