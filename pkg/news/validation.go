package news

/*
	This code block contains all Validation Code used for the entity.
	This code is mainly used in the services to check before sending to db or cache

*/
import (
	"errors"
	"fmt"
)

func validateID(id uint) error {
	if id == 0 {
		err := errors.New("News ID needs to be not null or zero.")
		return err
	}
	return nil
}
func validateString(s string, propertyName string) error {
	if len(s) <= 0 {
		err := errors.New(fmt.Sprintf("%v cannot be empty.", propertyName))
		return err
	}
	return nil
}
