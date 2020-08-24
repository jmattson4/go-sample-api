package domain

import (
	"github.com/jinzhu/gorm"
)

//TokenDetails models the Tokens used for authentication and Authorization
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

//Account a struct to rep user account
type Account struct {
	gorm.Model
	Email        string `json:"email"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	AccessToken  string `json:"accessToken" sql:"-"`
	RefreshToken string `json:"refreshToken" sql:"-"`
}

func AccountBasicConstructor() *Account {
	return &Account{}
}
func AccountConstructor(email string, password string, role string) *Account {
	return &Account{
		Email:    email,
		Password: password,
		Role:     role,
	}
}

type AccountDBRepo interface {
	GetAccount(u uint) *Account
	Create(email string, password string) error
	Construct(interface{}) interface{}
	GetAccountByEmail(email string) *Account
}

type AccountCacheRepo interface {
	GetAccount(u uint) *Account
	Create(email string, password string) error
	Construct(interface{}) interface{}
	CreateToken(accountID uint) (*TokenDetails, error)
	CreateAuth(userid uint, td *TokenDetails) error
	DeleteAuth(givenUuid string) (int64, error)
}

type AccountService interface {
	Login(email, password string) (*Account, error)
	Validate(email, password string) error
	Create(email, password string) error
}
