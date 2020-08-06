package model

import (
	"fmt"
	"log"
	"os"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB
var userDB *gorm.DB //database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	dbHost := os.Getenv("INSTANCE_CONNECTION_NAME")

	accountUsername := os.Getenv("ACCOUNT_USER")
	accountPassword := os.Getenv("ACCOUNT_PASSWORD")
	accountDBName := os.Getenv("ACCOUNT_DATABASE_NAME")

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
	db.Debug().AutoMigrate(&Account{}, &Product{}) //Database migration
}

//GetDB returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}

//GetUserDB returns handle to the UserDB object
func GetUserDB() *gorm.DB {
	return userDB
}
