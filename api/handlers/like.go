package handlers

import (
	"github.com/jinzhu/gorm"
	"net/http"
)

const (
	couldNotretrieveArticleLikes = "could not retrieve article likes"
	couldNotRetrieveCommentLikes = "could not retrieve comment likes"
	likeCommentFailed = "could not like comment"
	likeArticleFailed = "could not like article"
	removeArticleLikeFailed = "remove like failed from article"
)

func (h *Handler) AddOrDislikeAnArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := isUserIDValid(w, r)
		if !ok {
			return
		}


		article, dao, err, errored := initializeArticleParamHandler(w, r, h)
		if errored {
			return
		}
		defer dao.DB.Close()

		_, err = dao.GetLikeForArticle(userID, article.ID)
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				http.Error(w, couldNotretrieveArticleLikes, http.StatusInternalServerError)
				return
			}

			_, err := dao.AddLikeToArticle(userID, *article)
			if err != nil {
				http.Error(w, likeArticleFailed, http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}

		err = dao.RemoveLikeFromArticle(userID, *article)
		if err != nil {
			http.Error(w, removeArticleLikeFailed, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) AddorDislikeAComment() http.HandlerFunc {
	 return func(w http.ResponseWriter, r *http.Request) {
		 userID, ok := isUserIDValid(w, r)
		 if !ok {
			 return
		 }


		 comment, dao, errored := initializeCommentParamHandler(w, r, h)
		 if errored {
			 return
		 }
		 defer dao.DB.Close()

		 _, err := dao.GetLikeForComment(userID, comment.ID)
		 if err != nil {
			 if !gorm.IsRecordNotFoundError(err) {
				 http.Error(w, couldNotRetrieveCommentLikes, http.StatusInternalServerError)
				 return
			 }

			 _, err := dao.AddLikeToComment(userID, *comment)
			 if err != nil {
				 http.Error(w, likeCommentFailed, http.StatusInternalServerError)
				 return
			 }

			 w.WriteHeader(http.StatusNoContent)
			 return
		 }

		 err = dao.RemoveLikeFromComment(userID, *comment)
		 if err != nil {
			 http.Error(w, removeArticleLikeFailed, http.StatusInternalServerError)
			 return
		 }

		 w.WriteHeader(http.StatusNoContent)
	 }
}