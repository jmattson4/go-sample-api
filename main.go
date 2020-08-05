// main.go

package main

import (
	"log"

	a "github.com/jmattson4/go-sample-api/app"
	db "github.com/jmattson4/go-sample-api/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	a := a.App{}
	defer db.GetDB().Close()
	defer db.GetUserDB().Close()
	a.Initialize()

	a.Run(":8010")
}
