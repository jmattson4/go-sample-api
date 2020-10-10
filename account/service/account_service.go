package service

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmattson4/go-sample-api/domain"
	"github.com/jmattson4/go-sample-api/util"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

//AccountService ...
//Holds the cache and db repo for Accounts
type AccountService struct {
	dbRepo    domain.AccountDBRepo
	cacheRepo domain.AccountCacheRepo
}

//ConstructAccountService ...
func ConstructAccountService(repoDB domain.AccountDBRepo, repoCache domain.AccountCacheRepo) *AccountService {
	return &AccountService{
		dbRepo:    repoDB,
		cacheRepo: repoCache,
	}
}

//GetAccount ...
//Gets an account given a UUID associated to that account.
func (as *AccountService) GetAccount(uuid *uuid.UUID) (*domain.Account, error) {
	acc, err := as.dbRepo.GetAccount(uuid)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

//GetAccounts ...
//Gets all accounts
func (as *AccountService) GetAccounts() ([]domain.Account, error) {
	acc, err := as.dbRepo.GetAccounts()
	if err != nil {
		return nil, err
	}
	return acc, nil
}

//GetAccountByEmail Grabs an account by the given email name.
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

//Create 's a account in the database if user given info passes the validation requirements.
func (as *AccountService) Create(email, password string) error {
	valErr := as.validate(email, password)
	if valErr != nil {
		return valErr
	}
	// checkErr := as.checkIfEmailExists(email)
	// if checkErr != nil && checkErr != domain.ACCOUNT_EMAIL_CANT_FIND {
	// 	return checkErr
	// }
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
	ts, err := as.CreateToken(account.ID.String())
	if err != nil {
		return nil, domain.ACCOUNT_TOKEN_CREATION_ERROR
	}
	saveErr := as.cacheRepo.CreateAuth(account.ID.String(), ts)

	if saveErr != nil {
		return nil, domain.ACCOUNT_CACHE_AUTH_CREATION
	}

	account.AccessToken = ts.AccessToken //Store the access token in the response
	account.RefreshToken = ts.RefreshToken

	return account, nil

}

//Logout used to delete a given accessUUID so that the user is logged out and cannot use that access token any longer
func (as *AccountService) Logout(accessUUID string) error {
	_, err := as.cacheRepo.DeleteAuth(accessUUID)
	if err != nil {
		return err
	}
	return nil
}

//CreateToken : Function used to generate a access token that expires in 15 minutes and a Refresh Token that expries in 7 days
func (as *AccountService) CreateToken(accountUUID string) (*domain.TokenDetails, error) {
	accessSecret := util.GetEnv().AccessSecret
	refreshSecret := util.GetEnv().RefreshSecret
	//Create token details setting
	td := &domain.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.NewV4().String()

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["account_id"] = accountUUID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), atClaims)
	td.AccessToken, err = at.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["account_id"] = accountUUID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, err
	}

	return td, nil

}

//CreateAuth ...
// Used to create an authorization token and store it in the token cache
func (as *AccountService) CreateAuth(accountUUID string, td *domain.TokenDetails) error {
	err := as.cacheRepo.CreateAuth(accountUUID, td)
	if err != nil {
		return err
	}
	return nil
}

//GetAuth Used to get a userID for a given accessUUID as a string
func (as *AccountService) GetAuth(accessUUID string) (string, error) {
	userID, err := as.cacheRepo.GetAuth(accessUUID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

//DeleteAuth ...
//USed to delete an auth token from the cache repo
func (as *AccountService) DeleteAuth(givenuuid string) error {
	deleted, err := as.cacheRepo.DeleteAuth(givenuuid)
	if err != nil || deleted == 0 {
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

	return nil
}

//checkIfEmailExists checks the db seeing if email address already exists
func (as *AccountService) checkIfEmailExists(email string) error {
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
