package functions

import (
	"encoding/json"
	"fmt"
	"net/http"

	utils "github.com/advokrat/utilities"
	_ "github.com/mattn/go-sqlite3"
)

func GetCoinsFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getcoins" {
		resp := &serverResponse{
			Message: "404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
	}
	tokenFromUser := c.Value
	rollno, _, _ := utils.GetMetadata(tokenFromUser)
	w.Header().Set("Content-Type", "application/json")

	resp := &serverResponse{
		Message: "",
	}

	switch r.Method {

	case "GET":

		coins, err := utils.GetCoinsFromRollNo(rollno)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Fprintf(w, " -User not found")
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.Message = "Your coins are " + fmt.Sprintf("%f", coins)
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)

		resp.Message = "Sorry, only GET requests are supported"
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}

}
