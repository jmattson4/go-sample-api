// main.go

package main

import (
	"log"

	a "github.com/jmattson4/go-sample-api/app"
	c "github.com/jmattson4/go-sample-api/cache"
	m "github.com/jmattson4/go-sample-api/model"
	sec "github.com/jmattson4/go-sample-api/security"
)

func main() {
	defer m.GetDB().Close()
	defer m.GetUserDB().Close()
	defer c.Client.Close()
	enforcer, err := sec.InitAuthorizationEnforcer()
	if err != nil {
		log.Fatal("Error Getting Auth Enforcer: %v", err)
	} else {
		a := a.App{}

		a.Initialize(enforcer)

		a.Run(":8010")
	}
}
