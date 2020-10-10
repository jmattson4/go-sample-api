package domain

import (
	"errors"
)

var (
	ACCOUNT_CREATION_SUCCESS     string = "Account has been successfully created. Please continue to login."
	ACCOUNT_ID_INVALID           error  = errors.New("Account ID invalid.")
	ACCOUNT_EMAIL_INVALID        error  = errors.New("Email address is invalid.")
	ACCOUNT_EMAIL_EMPTY          error  = errors.New("Email needs to be not empty!")
	ACCOUNT_PASSWORD_TOO_SHORT   error  = errors.New("Password is less than six characters!")
	ACCOUNT_EMAIL_IN_USE         error  = errors.New("Email address is already in use!")
	ACCOUNT_EMAIL_CANT_FIND      error  = errors.New("Can't find Account with the following Email.")
	ACCOUNT_ID_CANT_FIND         error  = errors.New("Cant find the Account for given UUID.")
	ACCOUNT_PASSWORD_NO_MATCH    error  = errors.New("Invalid login credentials. Please try again")
	ACCOUNT_CREATION_FAILURE     error  = errors.New("Account creation has failed. Please try again.")
	ACCOUNT_TOKEN_CREATION_ERROR error  = errors.New("Error while creating token.")
	ACCOUNT_CACHE_AUTH_CREATION  error  = errors.New("Cannot create Auth in cache.")
	ACCOUNT_CACHE_AUTH_DELETION  error  = errors.New("Cannot delete Auth. Given id may not be valid.")

	NEWS_STANDARD_TESTING_ERROR                          error  = errors.New("Testing Error")
	NEWS_CACHE_CREATION_ERROR                            error  = errors.New("Failed to save to cache.")
	NEWS_DB_CREATION_ERROR                               error  = errors.New("Failed to save to database.")
	NEWS_CACHE_UPDATE_ERROR                              error  = errors.New("Failed to update to cache.")
	NEWS_DB_UPDATE_ERROR                                 error  = errors.New("Failed to update to database.")
	NEWS_CACHE_DELETE_ERROR                              error  = errors.New("Failed to delete to cache.")
	NEWS_DB_DELETE_ERROR                                 error  = errors.New("Failed to delete to database.")
	NEWS_SERVICE_GETMULTIPLENEWS_START                   error  = errors.New("Start must be greater than zero")
	NEWS_SERVICE_GETMULTIPLENEWS_COUNT                   error  = errors.New("Count must be greater than zero")
	NEWS_SERVICE_GETMULTIPLENEWS_EMPTYSLICE              error  = errors.New("Length of News Slice must be greater than 0")
	NEWS_SERVICE_VALIDATEID_UUID_ISNIL                   error  = errors.New("News ID cannot be null or empty")
	NEWS_CONTROLLER_GETNEWSBYWEBNAME_INVALID_COUNT_FIELD string = "Invalid or empty form value: count"
	NEWS_CONTROLLER_GETNEWSBYWEBNAME_INVALID_START_FIELD string = "Invalid or empty form value: start"
	NEWS_CONTROLLER_GETNEWSBYWEBNAME_INVALID_GET_ERROR   string = "Error Getting News Articles: please try again."
	NEWS_CONTROLLER_GETNEWSARTICLEBYID_UUID_PARSE        string = "Given UUID is not in proper UUID format."
	NEWS_CONTROLLER_GETNEWSARTICLEBYID_ARTICLE_NOT_FOUND string = "Could not find Article for given Webname and ID"
	NEWS_CONTROLLER_CHECKWEBNAME_INVALID_WEBNAME         string = "Invalid News Website Name"
	NEWS_CONTROLLER_PULLNEWSDATA_FAILED_SCRAPE           string = "Scrape failed"
	NEWS_CONTROLLER_PULLNEWSDATA_FAILED_PROCESS          string = "Processing data into database failed."
	NEWS_CONTROLLER_PULLNEWSDATA_SUCCESS                 string = "News Scrape and Process has succeeded."

	JWT_CANNOT_FIND_PROPERTY_ACCESS  error = errors.New("Cannot find Access UUID property. May be null or empty.")
	JWT_CANNOT_FIND_PROPERTY_ACCOUNT error = errors.New("Cannot find AccountID property. May be null or empty.")

	SCRAPER_GLOBAL_NEWS_CANT_FIND_URL error = errors.New("Cannot find any pages associated with the given URL : Global News Scraper")
)
