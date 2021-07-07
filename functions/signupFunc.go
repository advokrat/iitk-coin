package functions

import (
	"encoding/json"
	"net/http"

	//Func "iitk-coin/functions"
	Utils "github.com/advokrat/utilities"

	_ "github.com/mattn/go-sqlite3"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string `json:"name"`
	Rollno   string `json:"rollno"`
	Password string `json:"password"`
	Account_type string `json:"account_type"`
}

type serverResponse struct {
	Message string `json:"message"`
}

//Processes SignUP
func SignupFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		resp := &serverResponse{
			Message: "Error 404: Page Not Found!!",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}

	switch r.Method {

	case "POST":

		var user User
		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		name := user.Name
		rollno := user.Rollno
		accountType := user.Account_type
		password := user.Password
		if rollno == "" || password == "" || accountType == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp := &serverResponse{
				Message: "Password and ROll No. are Mandatory!!",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			//log.Fatal(err)
			w.WriteHeader(401)
			resp := &serverResponse{
				Message: "Server Error Encountered!!",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
		}

		write_err := Utils.WriteUserToDb(name, rollno, string(hashed_password), accountType)

		if write_err != nil {
			w.WriteHeader(500) // Return 500 Internal Server Error.
			resp := &serverResponse{
				Message: "This ROll No. is Already Registered. Try Logging In instead!!",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		w.WriteHeader(http.StatusOK)
		//Write json response back to response
		resp := &serverResponse{
			Message: "Account Creation Successful!!",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		resp := &serverResponse{
			Message: "Invalid Request Type!! Only POST Requests Accepted!!",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}

}
