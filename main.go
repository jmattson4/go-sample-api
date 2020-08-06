// main.go

package main

import (
	"log"

	a "github.com/jmattson4/go-sample-api/app"
	m "github.com/jmattson4/go-sample-api/model"
	sec "github.com/jmattson4/go-sample-api/security"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	defer m.GetDB().Close()
	defer m.GetUserDB().Close()
	enforcer, err := sec.InitAuthorizationEnforcer()
	if err != nil {
		log.Fatal("Error Getting Auth Enforcer: %v", err)
	} else {
		a := a.App{}

		a.Initialize(enforcer)

		a.Run(":8010")
	}
}
