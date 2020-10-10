package database

import (
	"fmt"
	"time"

	"github.com/jmattson4/go-sample-api/domain"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

//InitNewsDB Used to initiliaze a connection to the news db used for user accounts
func InitNewsDB(dbUser string, dbPW string, dbN string, dbService string, dbP string) *gorm.DB {
	username := dbUser
	password := dbPW
	dbName := dbN
	dbHost := dbService
	dbPort := dbP

	dbURI := fmt.Sprintf("host=%s sslmode=disable port=%s user=%s dbname=%s password=%s ", dbHost, dbPort, username, dbName, password)

	for i := 0; i < 10; i++ {
		conn, err := gorm.Open("postgres", dbURI)
		if err != nil {
			fmt.Printf("Attempt %s : Unable to open DB: %s ... Retrying \n", fmt.Sprint(i), err)
		} else {
			conn.DB().SetConnMaxLifetime(20 * time.Second)
			conn.DB().SetMaxIdleConns(30)
			fmt.Println("Connection to News database was succesful.")
			conn.Debug().AutoMigrate(&domain.NewsData{}) //Database migration
			return conn
		}
	}
	return nil
}

//InitAccountDB Used to initilaze a connection to the user db used for user accounts
func InitAccountDB(accDBServ string, accDBP string, accDBUSer string, accPW string, accDBN string) *gorm.DB {
	accountHost := accDBServ
	accountPort := accDBP
	accountUsername := accDBUSer
	accountPassword := accPW
	accountDBName := accDBN

	dbURI2 := fmt.Sprintf("host=%s sslmode=disable port=%s user=%s dbname=%s password=%s ", accountHost, accountPort, accountUsername, accountDBName, accountPassword)

	for i := 0; i < 10; i++ {
		userConn, userErr := gorm.Open("postgres", dbURI2)
		if userErr != nil {
			fmt.Printf("Attempt #%s: Unable to open User DB: %s ... Retrying \n", fmt.Sprint(i), userErr)
		} else {
			userConn.DB().SetConnMaxLifetime(20 * time.Second)
			userConn.DB().SetMaxIdleConns(30)
			fmt.Println("Connection to Users database was succesful.")
			userConn.Debug().AutoMigrate(&domain.Account{})
			return userConn
		}
	}
	return nil
}
