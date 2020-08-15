package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/twinj/uuid"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/jmattson4/go-sample-api/cache"
	u "github.com/jmattson4/go-sample-api/util"
	"golang.org/x/crypto/bcrypt"
)

/*
Token ... JWT claims struct
*/
type Token struct {
	UserID uint
	jwt.StandardClaims
}

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

//CheckIfTablesExist ...
func CheckIfTablesExist() {
	accountExists := GetUserDB().HasTable("accounts")
	if !accountExists {
		GetUserDB().CreateTable(Account{})
	}
}

//Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := GetUserDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

//Create ...
func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	account.Role = "user"

	err := GetUserDB().Create(account).Error

	if err != nil {
		return u.Message(false, err.Error())
	}
	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

//Login ... allows the user to login
func Login(email, password string) map[string]interface{} {

	account := &Account{}
	err := GetUserDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT tokens
	ts, err := CreateToken(account.ID)
	if err != nil {
		return u.Message(false, fmt.Sprintf("Error: %v", err.Error()))
	}
	saveErr := CreateAuth(account.ID, ts)
	//TODO change to vague error message when in prod
	if saveErr != nil {
		return u.Message(false, fmt.Sprintf("Error: %v", saveErr.Error()))
	}

	resp := u.Message(true, "Logged In")
	account.AccessToken = ts.AccessToken //Store the access token in the response
	account.RefreshToken = ts.RefreshToken
	resp["account"] = account

	return resp

}

//CreateToken : Function used to generate a access token that expires in 15 minutes and a Refresh Token that expries in 7 days
func CreateToken(accountID uint) (*TokenDetails, error) {
	accessSecret := u.GetEnv().AccessSecret
	refreshSecret := u.GetEnv().RefreshSecret
	//Create token details setting
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["account_id"] = accountID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), atClaims)
	td.AccessToken, err = at.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["account_id"] = accountID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, err
	}

	return td, nil

}

func CreateAuth(userid uint, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := cache.Client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := cache.Client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

//GetUser ... Gets the user from the userID
func GetUser(u uint) *Account {
	acc := &Account{}
	GetUserDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}
	acc.Password = ""
	return acc
}
