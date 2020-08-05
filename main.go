// main.go

package main

import (
	"log"
	"os"

	a "github.com/jmattson4/go-sample-api/app"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	a := a.App{}
	a.Initialize(
		os.Getenv("DATABASE_USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("INSTANCE_CONNECTION_NAME"))

	a.Run(":8010")
}
