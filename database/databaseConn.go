package database

import (
	"fmt"
	"time"

	"github.com/jmattson4/go-sample-api/model"
	"github.com/jmattson4/go-sample-api/util"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *gorm.DB
var userDB *gorm.DB //database

func init() {
	env := util.GetEnv()

	go initNewsDB(env)
	go initUserDb(env)

}

func initNewsDB(env *util.Environmentals) {
	username := env.DatabaseUser
	password := env.DatabasePassword
	dbName := env.DatabaseName
	dbHost := env.DatabaseDBService
	dbPort := env.DatabaseDBPort

	dbURI := fmt.Sprintf("host=%s sslmode=disable port=%s user=%s dbname=%s password=%s ", dbHost, dbPort, username, dbName, password)

	for i := 0; i < 10; i++ {
		conn, err := gorm.Open("postgres", dbURI)
		if err != nil {
			fmt.Printf("Attempt %s : Unable to open DB: %s ... Retrying \n", i, err)
			time.Sleep(time.Second * 5)
		} else {
			conn.DB().SetConnMaxLifetime(20 * time.Second)
			conn.DB().SetConnMaxIdleConn(30)
			db = conn
			fmt.Println("Connection to News database wassuccesful.")
			db.Debug().AutoMigrate(&model.Product{}, &model.NewsData{}) //Database migration
			break
		}
	}
}

func initUserDb(env *util.Environmentals) {
	accountHost := env.AccountDBService
	accountPort := env.AccountDBPort
	accountUsername := env.AccountUser
	accountPassword := env.AccountPassword
	accountDBName := env.AccountDatabaseName

	dbURI2 := fmt.Sprintf("host=%s sslmode=disable port=%s user=%s dbname=%s password=%s ", accountHost, accountPort, accountUsername, accountDBName, accountPassword)

	for i := 0; i < 10; i++ {
		userConn, userErr := gorm.Open("postgres", dbURI2)
		if userErr != nil {
			fmt.Printf("Attempt #%s: Unable to open User DB: %s ... Retrying \n", i, userErr)
			time.Sleep(time.Second * 5)
		} else {
			userConn.DB().SetConnMaxLifetime(20 * time.Second)
			userConn.DB().SetConnMaxIdleConn(30)
			userDB = userConn
			fmt.Println("Connection to Users database was succesful.")
			userDB.Debug().AutoMigrate(&model.Account{})
			break
		}
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
