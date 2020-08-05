package util

import (
	"encoding/json"
	"net/http"
)

//Message ... This responds to the user with status and a response.
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

//RespondWithError This is a util to respond to the client with a error message
//	if something went wrong
func RespondWithError(w http.ResponseWriter, code int, message map[string]interface{}) {
	RespondWithJSON(w, code, message)
}

//RespondWithJSON This is a function that takes a piece of payload and then turns it into a Json object and
//	sends back to the client.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
