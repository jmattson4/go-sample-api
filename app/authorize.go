package app

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/jmattson4/go-sample-api/model"
	u "github.com/jmattson4/go-sample-api/util"
)

//Authorize ...
func Authorize(e *casbin.Enforcer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value("user")
			if user != nil {
				userID := user.(uint)
				userAcc := model.GetUser(userID)
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
