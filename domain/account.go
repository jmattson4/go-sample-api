package domain

import (
	"github.com/jinzhu/gorm"
	"github.com/twinj/uuid"
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
	ID           uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Role         string    `json:"role"`
	AccessToken  string    `json:"accessToken" sql:"-"`
	RefreshToken string    `json:"refreshToken" sql:"-"`
}

func AccountBasicConstructor() *Account {
	return &Account{}
}

func AccountConstructorWithID() *Account {
	uuid := uuid.NewV4()
	return &Account{
		ID: uuid,
	}
}
func AccountConstructor(email string, password string, role string) *Account {
	uuid := uuid.NewV4()
	return &Account{
		ID:       uuid,
		Email:    email,
		Password: password,
		Role:     role,
	}
}

type AccountDBRepo interface {
	GetAccount(uuid *uuid.UUID) (*Account, error)
	Create(email string, password string) error
	Construct(interface{}) interface{}
	GetAccountByEmail(email string) (*Account, error)
}

type AccountCacheRepo interface {
	GetAuth(accessUUID string) (string, error)
	CreateAuth(accountUUID string, td *TokenDetails) error
	DeleteAuth(givenUuid string) (int64, error)
}

type AccountService interface {
	GetAccount(uuid *uuid.UUID) (*Account, error)
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
