package accounts

import (
	"github.com/jinzhu/gorm"
)

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
	AccessToken  string `json:"accessToken";sql:"-"`
	RefreshToken string `json:"refreshToken";sql:"-"`
}
