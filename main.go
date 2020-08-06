// main.go

package main

import (
	"log"

	a "github.com/jmattson4/go-sample-api/app"
	m "github.com/jmattson4/go-sample-api/model"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	a := a.App{}
	defer m.GetDB().Close()
	defer m.GetUserDB().Close()
	a.Initialize()

	a.Run(":8010")
}
