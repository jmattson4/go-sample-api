package model

import (
	"fmt"
	"time"

	"github.com/jmattson4/go-sample-api/util"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *gorm.DB
var userDB *gorm.DB //database

func init() {
	env := util.GetEnv()

	username := env.DatabaseUser
	password := env.DatabasePassword
	dbName := env.DatabaseName
	dbHost := env.DatabaseDBService
	dbPort := env.DatabaseDBPort

	accountHost := env.AccountDBService
	accountPort := env.AccountDBPort
	accountUsername := env.AccountUser
	accountPassword := env.AccountPassword
	accountDBName := env.AccountDatabaseName

	dbURI := fmt.Sprintf("host=%s sslmode=disable port=%s user=%s dbname=%s password=%s ", dbHost, dbPort, username, dbName, password)
	dbURI2 := fmt.Sprintf("host=%s sslmode=disable port=%s user=%s dbname=%s password=%s ", accountHost, accountPort, accountUsername, accountDBName, accountPassword)
	time.Sleep(time.Second * 30)
	for i := 0; i < 10; i++ {
		conn, err := gorm.Open("postgres", dbURI)
		if err != nil {
			fmt.Printf("Unable to open DB: %s ... Retrying \n", err)
			time.Sleep(time.Second * 5)
		}
		userConn, userErr := gorm.Open("postgres", dbURI2)
		if userErr != nil {
			fmt.Printf("Unable to open User DB: %s ... Retrying \n", userErr)
			time.Sleep(time.Second * 5)
		}
		if err == nil {
			fmt.Println("Connection to DB succesful.")
			db = conn
		}
		if userErr == nil {
			fmt.Println("Connection to UserDB succesful.")
			userDB = userConn
		}
		if err == nil && userErr == nil {
			db = conn
			userDB = userConn
			break
		}
	}
	if db != nil && userDB != nil {
		fmt.Println("Connection to database succesful.")
		db.Debug().AutoMigrate(&Product{}, &NewsData{}) //Database migration
		userDB.Debug().AutoMigrate(&Account{})
	} else if db != nil {
		db.Debug().AutoMigrate(&Product{}, &NewsData{}) //Database migration
	} else if userDB != nil {
		userDB.Debug().AutoMigrate(&Account{})
	}
}

//GetDB returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}

//GetUserDB returns handle to the UserDB object
func GetUserDB() *gorm.DB {
	return userDB
}
