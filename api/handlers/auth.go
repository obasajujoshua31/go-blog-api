package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"go-blog-api/config"
	"go-blog-api/dal"
	"go-blog-api/services"
	"io/ioutil"
	"net/http"
)

type ErrorMessage struct {
	Errors []string `json:"errors"`
}

type AuthPayload struct {
	UserToken string `json:"user_token"`
}

const (
	unableToUnMarshalRequest = "Unable to unmarshal request"
	unableToMarshalResponse  = "Unable to marshal response"
	unableToConnectToDB      = "Unable to connect to Database"
	userlookupfailed         = "user look up failed"
	emailAlreadExist         = "email already exist"
	unableToCreateUser       = "unable to create user"
	unableToGenerateJWT      = "unable to generate jwt token"
	loginCredentialsError    = "invalid login credentials"
)

type Handler struct {
	AppConfig *config.AppConfig
}

func (h *Handler) GetHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", "Welcome to go blog api ...")
	}
}

func (h *Handler) RegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, dao, failed := initializeAuthHandler(r, w, h, isRegisterRequestValid)
		if failed {
			return
		}

		defer dao.DB.Close()

		foundUser, err := dao.GetUserByEmail(user.Email)
		if foundUser != nil {
			http.Error(w, emailAlreadExist, http.StatusBadRequest)
			return
		}

		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				http.Error(w, userlookupfailed, http.StatusInternalServerError)
				return
			}
		}

		user.ID = services.GenerateUUID()

		createdUser, err := dao.CreateNewUser(*user)
		if err != nil {
			http.Error(w, unableToCreateUser, http.StatusInternalServerError)
			return
		}

		generateTokenAndRespond(createdUser, w)
	}
}

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, dao, failed := initializeAuthHandler(r, w, h, isLoginValid)
		if failed {
			return
		}
		defer dao.DB.Close()

		foundUser, err := dao.GetUserByEmail(user.Email)
		if err != nil || foundUser == nil {
			if gorm.IsRecordNotFoundError(err) {
				http.Error(w, loginCredentialsError, http.StatusBadRequest)
				return
			}

			http.Error(w, userlookupfailed, http.StatusInternalServerError)
			return

		}

		if !services.IsPasswordMatch(user.Password, foundUser.Password) {
			http.Error(w, loginCredentialsError, http.StatusBadRequest)
			return
		}

		generateTokenAndRespond(foundUser, w)
	}
}

func initializeAuthHandler(r *http.Request, w http.ResponseWriter, h *Handler, isRequestValid func(user dal.User) []error) (*dal.User, *dal.DAL, bool) {
	user := dal.User{}

	err := unMarshalRequest(r, &user)
	if err != nil {
		http.Error(w, unableToUnMarshalRequest, http.StatusInternalServerError)
		return nil, nil, true
	}

	valErrors := isRequestValid(user)
	if len(valErrors) > 0 {
		writeValErrors(w, valErrors)
		return nil, nil, true
	}

	dao, failed := GetDAOFromDB(h, w)
	if failed {
		return nil, nil, true
	}
	return &user, dao, false
}

func unMarshalRequest(r *http.Request, out interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		return err
	}
	return nil
}

func writeValErrors(w http.ResponseWriter, valErrors []error) {
	errorArray := []string{}
	for _, valError := range valErrors {
		errorArray = append(errorArray, valError.Error())
	}
	errResp := ErrorMessage{
		Errors: errorArray,
	}

	w.WriteHeader(http.StatusBadRequest)
	writeResponse(w, errResp)
}

func writeResponse(w http.ResponseWriter, data interface{}) {
	resp, err := json.MarshalIndent(&data, "", " ")
	if err != nil {
		http.Error(w, unableToMarshalResponse, http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func generateTokenAndRespond(user *dal.User, w http.ResponseWriter) {
	userToken, err := services.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, unableToGenerateJWT, http.StatusInternalServerError)
		return
	}

	userResp := AuthPayload{
		UserToken: userToken,
	}

	writeResponse(w, userResp)
}
