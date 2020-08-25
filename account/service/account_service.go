package service

import (
	"strings"

	"github.com/jmattson4/go-sample-api/domain"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	dbRepo    domain.AccountDBRepo
	cacheRepo domain.AccountCacheRepo
}

func (as *AccountService) GetAccount(uuid *uuid.UUID) (*domain.Account, error) {
	acc, err := as.dbRepo.GetAccount(uuid)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (as *AccountService) GetAccountByEmail(email string) (*domain.Account, error) {
	if len(email) < 0 {
		return nil, domain.ACCOUNT_EMAIL_EMPTY
	}
	if !strings.Contains(email, "@") {
		return nil, domain.ACCOUNT_EMAIL_INVALID
	}
	acc, err := as.dbRepo.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (as *AccountService) Create(email, password string) error {
	valErr := as.validate(email, password)
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

	account, err := as.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	hashErr := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if hashErr != nil && hashErr == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
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

	if saveErr != nil {
		return nil, domain.ACCOUNT_CACHE_AUTH_CREATION
	}

	account.AccessToken = ts.AccessToken //Store the access token in the response
	account.RefreshToken = ts.RefreshToken

	return account, nil

}

func (as *AccountService) Logout(accessUuid string) error {
	_, err := as.cacheRepo.DeleteAuth(accessUuid)
	if err != nil {
		return err
	}
	return nil
}

//Validate incoming user details...
func (as *AccountService) validate(email, password string) error {
	if !strings.Contains(email, "@") {
		return domain.ACCOUNT_EMAIL_INVALID
	}

	if len(password) < 6 {
		return domain.ACCOUNT_PASSWORD_TOO_SHORT
	}

	//check for errors and duplicate emails
	temp, err := as.GetAccountByEmail(email)
	if err != nil {
		return err
	}
	if temp.Email != "" {
		return domain.ACCOUNT_EMAIL_IN_USE
	}

	return nil
}
