package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
)

//InitDatabase ...
func InitDatabase(user string, password string, dbname string, instanceConnectionName string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		instanceConnectionName,
		dbname,
		user,
		password)

	var err error
	var DB *sql.DB
	DB, err = sql.Open("cloudsqlpostgres", connectionString)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return DB, nil
}
