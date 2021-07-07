package functions

import (
	"encoding/json"
	"net/http"

	//Func "iitk-coin/functions"
	Utils "github.com/advokrat/utilities"

	_ "github.com/mattn/go-sqlite3"

	"golang.org/x/crypto/bcrypt"
)

//Processes Login
func LoginFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		resp := &serverResponse{
			Message: "Error 404: Page Not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {

	case "POST":

		var user User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rollno := user.Rollno
		password := user.Password
		hashedPW := Utils.HashPassword(rollno)

		// Authenticating Password with Hash
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(password)); err != nil {
			w.WriteHeader(500) //Configuring Server Error Message
			resp := &serverResponse{
				Message: "Password Is Incorrect!! Try Again!!",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		token, expirationTime, err := Utils.CreateJWT(rollno)

		if err != nil {
			w.WriteHeader(401)
			resp := &serverResponse{
				Message: "Server Error",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return

		}

		//Setting appropriate Cookie for the User
		http.SetCookie(w, &http.Cookie{
			Name:     "JWT Token",
			Value:    token,
			Expires:  expirationTime,
			HttpOnly: true,
		})

		w.WriteHeader(http.StatusOK)

		resp := &serverResponse{
			Message: "Login Successful!!",
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
