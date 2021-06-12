package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string `json:"name"`
	Rollno   string `json:"rollno"`
	Password string `json:"password"`
}

type serverResponse struct {
	Message string `json:"message"`
}

//Processes the request to Access Secret Page
func secretPageFunc(w http.ResponseWriter, r *http.Request) {
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
		user_roll_no, err := GetMetadata(tokenFromUser)

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
			Message: "Access Granted!! Welcome ROll No. " + user_roll_no,
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

//Processes Login
func loginFunc(w http.ResponseWriter, r *http.Request) {
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
		hashedPW := HashPassword(rollno)

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

		token, expirationTime, err := CreateJWT(rollno)

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

//Processes SignUP
func signupFunc(w http.ResponseWriter, r *http.Request) {
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
		password := user.Password
		if rollno == "" || password == "" {
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

		write_err := inputToDB(name, rollno, string(hashed_password))

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

//Hashes the Password before storing it in Database
func HashPassword(rollno string) string {
	database, _ :=
		sql.Open("sqlite3", "./user.db")
	rollno_int, _ := strconv.Atoi(rollno)

	row := database.QueryRow(`SELECT password FROM user WHERE rollno= $1;`, rollno_int)

	var hashed_password string
	row.Scan(&hashed_password)

	return (hashed_password)

}

//Inputs Data to Database upon Signup
func inputToDB(name string, rollno string, password string) error {
	database, _ :=
		sql.Open("sqlite3", "./user.db")

	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS user (name TEXT,rollno TEXT PRIMARY KEY,password TEXT)")

	statement.Exec()

	statement, _ =
		database.Prepare("INSERT INTO user (name,rollno,password) VALUES (?, ?, ?)")
	_, err := statement.Exec(name, rollno, password)
	if err != nil {
		return err
	}
	return nil

}

//Creates JWT Token upon Login
func CreateJWT(userRollNo string) (string, time.Time, error) {
	var err error
	//Creating JWT

	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatal("An Error Encountered While Loading .env File!!")
	}

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_roll_no"] = userRollNo
	expTime := time.Now().Add(time.Minute * 15)
	atClaims["exp"] = expTime.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("SecretKey")))
	if err != nil {
		return "", time.Now(), err
	}
	return token, expTime, err

}

//Authenticates JWT Token
func AuthJWT(request_token string) (*jwt.Token, error) {
	tokenString := request_token
	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatal("An Error Encountered While Loading .env File!!")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Conforming the Token Structure to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Error!! Unexpected Signing Method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SecretKey")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil

}

//Returns ROll NO. of the User associated with the Token
func GetMetadata(user_token string) (string, error) {
	token, err := AuthJWT(user_token)
	if err != nil {
		return " ", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		roll_no, _ := claims["user_roll_no"].(string)
		return roll_no, err
	}

	return " ", err

}

//Main Driver Function
func main() {

	http.HandleFunc("/secretpage", secretPageFunc)
	http.HandleFunc("/login", loginFunc)
	http.HandleFunc("/signup", signupFunc)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("An Error Encountered While Loading .env File!!")
	}

	fmt.Println("Listening on Port 8080!!")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
