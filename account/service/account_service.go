package service

import (
	"strings"

	"github.com/jmattson4/go-sample-api/domain"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	dbRepo    domain.AccountDBRepo
	cacheRepo domain.AccountCacheRepo
}

func (as *AccountService) Create(email, password string) error {
	valErr := as.Validate(email, password)
	if valErr != nil {
		return valErr
	}
	createErr := as.dbRepo.Create(email, password)
	if createErr != nil {
		return createErr
	}
	return nil
}

//Login ... allows the user to login
func (as *AccountService) Login(email, password string) (*domain.Account, error) {

	account := &domain.Account{}
	account = as.dbRepo.GetAccountByEmail(email)
	if account == nil {
		return nil, domain.ACCOUNT_EMAIL_CANT_FIND
	}

	err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!

		return nil, domain.ACCOUNT_PASSWORD_NO_MATCH
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT tokens
	ts, err := as.cacheRepo.CreateToken(account.ID)
	if err != nil {
		return nil, domain.ACCOUNT_TOKEN_CREATION_ERROR
	}
	saveErr := as.cacheRepo.CreateAuth(account.ID, ts)
	//TODO change to vague error message when in prod
	if saveErr != nil {
		return nil, domain.ACCOUNT_CACHE_AUTH_CREATION
	}

	account.AccessToken = ts.AccessToken //Store the access token in the response
	account.RefreshToken = ts.RefreshToken

	return account, nil

}

//Validate incoming user details...
func (as *AccountService) Validate(email, password string) error {
	if !strings.Contains(email, "@") {
		return domain.ACCOUNT_EMAIL_INVALID
	}

	if len(password) < 6 {
		return domain.ACCOUNT_PASSWORD_TOO_SHORT
	}

	//Email must be unique
	temp := &domain.Account{}

	//check for errors and duplicate emails
	temp = as.dbRepo.GetAccountByEmail(email)
	if temp == nil {
		return domain.ACCOUNT_EMAIL_CANT_FIND
	}
	if temp.Email != "" {
		return domain.ACCOUNT_EMAIL_IN_USE
	}

	return nil
}
