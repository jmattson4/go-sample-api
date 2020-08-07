package database

import (
	"log"

	"github.com/jinzhu/gorm"
)

//Transaction : This function is runs a transaction on the database.
//	it takes a database and then a funcSlice to run a series of functions
//	on the database and if any function fails it rolls it back. If all functions
//	work then it commits it to the database returning an error if it fails and nil
//	if it doesnt.
func Transaction(db *gorm.DB, funcSlice []func(interface{}) error) error {
	db.Begin()
	for _, function := range funcSlice {
		err := function(interface{})
		if err != nil {
			db.Rollback()
			log.Print("Error with inserting news into DB. Rolling back DB.")
			return err
		}
	}
	commitErr := db.Commit().Error
	if commitErr != nil {
		log.Fatal("Error Database failed to commit transaction!")
		return commitErr
	}
	return nil
}
