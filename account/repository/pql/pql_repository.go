package pql

import (
	"github.com/jinzhu/gorm"
	"github.com/jmattson4/go-sample-api/domain"
	"github.com/twinj/uuid"
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

//Create ...
func (repo *AccountsRepo) Create(email string, password string) error {

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	account := domain.AccountConstructor(email, string(hashedPassword), "user")

	err := repo.db.Create(account).Error

	if err != nil {
		return domain.ACCOUNT_CREATION_FAILURE
	}
	if uuid.IsNil(account.ID) {
		return domain.ACCOUNT_ID_INVALID
	}

	account.Password = "" //delete password

	return nil
}

//CheckIfTablesExist ...
func (repo *AccountsRepo) CheckIfTablesExist() {
	accountExists := repo.db.HasTable("accounts")
	if !accountExists {
		repo.db.CreateTable(domain.Account{})
	}
}

//GetAccount ... Gets the Account from the AccountID
func (repo *AccountsRepo) GetAccount(u *uuid.UUID) (*domain.Account, error) {
	acc := domain.AccountBasicConstructor()
	err := repo.db.Table("accounts").Where("id = ?", u).First(acc).Error
	if err != nil || acc.Email == "" { //User not found!
		return nil, domain.ACCOUNT_ID_CANT_FIND
	}
	acc.Password = ""
	return acc, nil
}

//GetAccountByEmail Returns the account associated with the given email address.
func (repo *AccountsRepo) GetAccountByEmail(email string) (*domain.Account, error) {
	acc := domain.AccountBasicConstructor()
	err := repo.db.Table("accounts").Where("email = ?", email).First(acc).Error
	if err != nil || acc.Email == "" { //User not found!
		return nil, domain.ACCOUNT_EMAIL_CANT_FIND
	}
	acc.Password = ""
	return acc, nil
}
