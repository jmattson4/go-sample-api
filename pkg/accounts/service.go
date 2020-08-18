package accounts

import (
	"fmt"

	"github.com/jinzhu/gorm"
	u "github.com/jmattson4/go-sample-api/util"
	"golang.org/x/crypto/bcrypt"
)

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
