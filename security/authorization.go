package security

import (
	cas "github.com/casbin/casbin/v2"
)

//InitAuthorizationEnforcer ... This function is used to create the casbin enforcer to
//	enforce authoriztion throughout the system
func InitAuthorizationEnforcer() (*cas.Enforcer, error) {
	e, err := cas.NewEnforcer("./rest_model.conf", "./policy.csv")
	if err != nil {
		return nil, err
	}
	return e, nil
}
