package handlers

import (
	"errors"
	"go-blog-api/dal"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isRegisterRequestValid(user dal.User) (validationError []error) {
	if user.Name == "" {
		validationError = append(validationError, errors.New("name cannot be empty"))
	}

	validationError = isEmailAndPasswordValid(user, validationError)

	return validationError
}

func isEmailAndPasswordValid(user dal.User, validationError []error) []error {
	if !emailRegex.MatchString(user.Email) {
		validationError = append(validationError, errors.New("email is not valid"))
	}

	if user.Password == "" {
		validationError = append(validationError, errors.New("password cannot be empty"))
	}
	return validationError
}

func isLoginValid(user dal.User) []error {
	return isEmailAndPasswordValid(user, []error{})
}

