package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//Environmentals used to model the environmentals of the application.
type Environmentals struct {
	AccountDBService    string
	AccountDBPort       string
	AccountDatabaseName string
	AccountUser         string
	AccountPassword     string
	DatabaseDBService   string
	DatabaseDBPort      string
	DatabaseName        string
	DatabaseUser        string
	DatabasePassword    string
	TokenPassword       string
}

var env *Environmentals

//init ...
//This function is used to inject a data structure which contains the environmental
//	variables into the system where needed.
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file. Trying Different route")
		err := godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading .env file. File probably isnt in system.")
		}
	}
	env = &Environmentals{
		AccountDBService:    os.Getenv("ACCOUNT_DB_SERVICE"),
		AccountDBPort:       os.Getenv("ACCOUNT_DB_PORT"),
		AccountDatabaseName: os.Getenv("ACCOUNT_DB"),
		AccountUser:         os.Getenv("ACCOUNT_DEV_USER"),
		AccountPassword:     os.Getenv("POSTGRES_DEV_PASSWORD"),
		DatabaseDBService:   os.Getenv("NEWS_DB_SERVICE"),
		DatabaseDBPort:      os.Getenv("NEWS_DB_PORT"),
		DatabaseName:        os.Getenv("NEWS_DB"),
		DatabaseUser:        os.Getenv("NEWS_DB_USER"),
		DatabasePassword:    os.Getenv("POSTGRES_DEV_PASSWORD"),
		TokenPassword:       os.Getenv("token_password"),
	}
}

//GetEnv ... used to grab the initilazed ENV
func GetEnv() *Environmentals {
	return env
}
