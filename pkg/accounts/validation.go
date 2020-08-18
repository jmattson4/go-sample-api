package accounts

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

//Validate incoming user details...
func Validate(account *Account, db *gorm.DB) error {
	if !strings.Contains(account.Email, "@") {
		err := errors.New("Email address is invalid.")
		return err

	}

	if len(account.Password) < 6 {
		err := errors.New("Password is too short!")
		return err
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := db.Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if temp.Email != "" {
		err := errors.New("Email address is already in use!")
		return err
	}

	return nil
}
