package app

import (
	"net/http"

	"github.com/twinj/uuid"

	"github.com/casbin/casbin/v2"
	acc "github.com/jmattson4/go-sample-api/account/service"
	u "github.com/jmattson4/go-sample-api/api/utils"
)

//Authorize ...
func Authorize(e *casbin.Enforcer, service *acc.AccountService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value("user")
			if user != nil {
				userID := user.(string)
				userUUID, err := uuid.Parse(userID)
				if err != nil {
					return
				}
				userAcc, getErr := service.GetAccount(userUUID)
				if getErr != nil {
					return
				}
				userRole := userAcc.Role

				// casbin rule enforcing

				enforceCasbin(e, next, w, r, userRole)
			} else {
				enforceCasbin(e, next, w, r, "anonymous")
			}
		}

		return http.HandlerFunc(fn)
	}
}

func enforceCasbin(e *casbin.Enforcer, next http.Handler, w http.ResponseWriter, r *http.Request, userRole string) {
	res, err := e.Enforce(userRole, r.URL.Path, r.Method)
	if err != nil {
		response := u.Message(false, err.Error())
		u.RespondWithError(w, http.StatusInternalServerError, response)
		return
	}
	if res {
		next.ServeHTTP(w, r)
	} else {
		response := u.Message(false, "Unauthorized")
		u.RespondWithError(w, http.StatusForbidden, response)
		return
	}
}
