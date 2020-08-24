package pql

import (
	"github.com/jinzhu/gorm"
	"github.com/jmattson4/go-sample-api/domain"
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
func (repo *AccountsRepo) GetAccountByEmail(email string) *domain.Account {
	acc := domain.AccountBasicConstructor()
	repo.db.Table("accounts").Where("email = ?", email).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}
	acc.Password = ""
	return acc
}

//Create ...
func (repo *AccountsRepo) Create(email string, password string) error {

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	account := domain.AccountConstructor(email, string(hashedPassword), "user")

	err := repo.db.Create(account).Error

	if err != nil {
		return err
	}
	if account.ID <= 0 {
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
func (repo *AccountsRepo) GetAccount(u uint) *domain.Account {
	acc := domain.AccountBasicConstructor()
	repo.db.Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}
	acc.Password = ""
	return acc
}
