package accounts

import (
	"errors"

	"github.com/jinzhu/gorm"
	u "github.com/jmattson4/go-sample-api/util"
	"golang.org/x/crypto/bcrypt"
)

type AccountsRepo struct {
	db *gorm.DB
}

func ConstructAccountsRepo(db *gorm.DB) *AccountsRepo {
	return &AccountsRepo{
		db: db,
	}
}
func (repo *AccountsRepo)GetByEmailName error {
	
}

//Create ...
func (repo *AccountsRepo) Create(account *Account) (*Account, error) {

	if valErr := Validate(account, repo.db); valErr != nil {
		return nil, valErr
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	account.Role = "user"

	err := repo.db.Create(account).Error

	if err != nil {
		return nil, err
	}
	if account.ID <= 0 {
		accountErr := errors.New("Account ID invalid.")
		return nil, accountErr
	}

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response, nil
}

//CheckIfTablesExist ...
func (repo *AccountsRepo) CheckIfTablesExist() {
	accountExists := repo.db.HasTable("accounts")
	if !accountExists {
		repo.db.CreateTable(Account{})
	}
}

//GetUser ... Gets the user from the userID
func (repo *AccountsRepo) GetUser(u uint) *Account {
	acc := &Account{}
	repo.db.Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}
	acc.Password = ""
	return acc
}
