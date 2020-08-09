package model

import (
	"fmt"

	"github.com/jmattson4/go-sample-api/util"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var userDB *gorm.DB //database

func init() {
	env := util.GetEnv()

	username := env.DatabaseUser
	password := env.DatabasePassword
	dbName := env.DatabaseName
	dbHost := env.InstanceConnectionName

	accountUsername := env.AccountUser
	accountPassword := env.AccountPassword
	accountDBName := env.AccountDatabaseName

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string
	db2URI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, accountUsername, accountDBName, accountPassword)

	conn, err := gorm.Open("cloudsqlpostgres", dbURI)
	if err != nil {
		fmt.Print(err)
	}
	userConn, userErr := gorm.Open("cloudsqlpostgres", db2URI)
	if userErr != nil {
		fmt.Print(userErr)
	}

	db = conn
	userDB = userConn

	db.Debug().AutoMigrate(&Account{}, &Product{}, &NewsData{}) //Database migration
}

//GetDB returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}

//GetUserDB returns handle to the UserDB object
func GetUserDB() *gorm.DB {
	return userDB
}
