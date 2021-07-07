package functions

import (
	"encoding/json"
	"net/http"

	//Func "iitk-coin/functions"
	Utils "github.com/advokrat/utilities"

	_ "github.com/mattn/go-sqlite3"
)

//Processes the request to Access Secret Page
func SecretPageFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/secretpage" {
		w.WriteHeader(404)
		resp := &serverResponse{
			Message: "Error 404: Page Not Found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				//Returns an Error Message, if JWT is not validated!
				w.WriteHeader(http.StatusUnauthorized)
				resp := &serverResponse{
					Message: "Authentication Failed!!",
				}
				JsonRes, _ := json.Marshal(resp)
				w.Write(JsonRes)

				return
			}
			//Returns default error message for any other errors!!
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tokenFromUser := c.Value
		user_roll_no, Acctype, err := Utils.GetMetadata(tokenFromUser)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			resp := &serverResponse{
				Message: "Current User is not Authorized to Access this Page!!",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}
		resp := &serverResponse{
			Message: "Access Granted!! Welcome ROll No. " + user_roll_no + " " + Acctype,
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		resp := &serverResponse{
			Message: "Invalid Request Type!! Only GET Requests Accepted!!",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
	}

}
