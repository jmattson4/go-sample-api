package domain

import (
	"errors"
)

var (
	ACCOUNT_CREATION_SUCCESS     string = "Account has been successfully created. Please continue to login."
	ACCOUNT_ID_INVALID           error  = errors.New("Account ID invalid.")
	ACCOUNT_EMAIL_INVALID        error  = errors.New("Email address is invalid.")
	ACCOUNT_PASSWORD_TOO_SHORT   error  = errors.New("Password is less than six characters!")
	ACCOUNT_EMAIL_IN_USE         error  = errors.New("Email address is already in use!")
	ACCOUNT_EMAIL_CANT_FIND      error  = errors.New("Can't find Account with the following Email.")
	ACCOUNT_PASSWORD_NO_MATCH    error  = errors.New("Invalid login credentials. Please try again")
	ACCOUNT_TOKEN_CREATION_ERROR error  = errors.New("Error while creating token.")
	ACCOUNT_CACHE_AUTH_CREATION  error  = errors.New("Cannot create Auth in cache.")

	NEWS_STANDARD_TESTING_ERROR error = errors.New("Testing Error")
	NEWS_CACHE_CREATION_ERROR   error = errors.New("Failed to save to cache.")
	NEWS_DB_CREATION_ERROR      error = errors.New("Failed to save to database.")
	NEWS_CACHE_UPDATE_ERROR     error = errors.New("Failed to update to cache.")
	NEWS_DB_UPDATE_ERROR        error = errors.New("Failed to update to database.")
	NEWS_CACHE_DELETE_ERROR     error = errors.New("Failed to delete to cache.")
	NEWS_DB_DELETE_ERROR        error = errors.New("Failed to delete to database.")
)
