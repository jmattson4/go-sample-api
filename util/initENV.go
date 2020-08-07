package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//Environmentals used to model the environmentals of the application.
type Environmentals struct {
	InstanceConnectionName string
	AccountDatabaseName    string
	AccountUser            string
	AccountPassword        string
	DatabaseName           string
	DatabaseUser           string
	DatabasePassword       string
	TokenPassword          string
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
		InstanceConnectionName: os.Getenv("INSTANCE_CONNECTION_NAME"),
		AccountDatabaseName:    os.Getenv("ACCOUNT_DATABASE_NAME"),
		AccountUser:            os.Getenv("ACCOUNT_USER"),
		AccountPassword:        os.Getenv("ACCOUNT_PASSWORD"),
		DatabaseName:           os.Getenv("DATABASE_NAME"),
		DatabaseUser:           os.Getenv("DATABASE_USER"),
		DatabasePassword:       os.Getenv("DATABASE_PASSWORD"),
		TokenPassword:          os.Getenv("token_password"),
	}
}

//GetEnv ... used to grab the initilazed ENV
func GetEnv() *Environmentals {
	return env
}
