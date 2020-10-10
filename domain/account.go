package domain

import (
	"github.com/twinj/uuid"
)

//TokenDetails models the Tokens used for authentication and Authorization
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

//Account a struct to rep user account
type Account struct {
	Base
	Email        string `json:"email"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	AccessToken  string `json:"accessToken" sql:"-"`
	RefreshToken string `json:"refreshToken" sql:"-"`
}

//AccountBasicConstructor ...
func AccountBasicConstructor() *Account {
	return &Account{}
}

//AccountConstructorWithID ...
func AccountConstructorWithID() *Account {
	base := Base{}
	return &Account{
		Base: base,
	}
}

//AccountConstructor ...
func AccountConstructor(email string, password string, role string) *Account {
	base := Base{}
	return &Account{
		Base:     base,
		Email:    email,
		Password: password,
		Role:     role,
	}
}

//AccountDBRepo ... interface used to represent a DB repo which can be
type AccountDBRepo interface {
	GetAccount(uuid *uuid.UUID) (*Account, error)
	Create(email string, password string) error
	GetAccountByEmail(email string) (*Account, error)
	GetAccounts() ([]Account, error)
}

//AccountCacheRepo interface used to represent a repo for the cache.
type AccountCacheRepo interface {
	GetAuth(accessUUID string) (string, error)
	CreateAuth(accountUUID string, td *TokenDetails) error
	DeleteAuth(givenUUID string) (int64, error)
}

//AccountService interface used to model a Account service which interacts with the repositories
type AccountService interface {
	GetAccount(uuid *uuid.UUID) (*Account, error)
	GetAccounts() ([]Account, error)
	GetAccountByEmail(email string) (*Account, error)
	Login(email, password string) (*Account, error)
	Logout(uuid string) error
	Validate(email, password string) error
	Create(email, password string) error
	GetAuth(accessUUID string) (string, error)
	CreateAuth(accountUUID string, td *TokenDetails) error
	DeleteAuth(givenUUID string) error
	CreateToken(accountUUID string) (*TokenDetails, error)
}
