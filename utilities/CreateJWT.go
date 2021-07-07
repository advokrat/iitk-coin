package utilities

import (
	"log"
	"os"
	"time"

	//Func "github.com/advokrat/functions"
	//Utils "iitk-coin/utilities"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"

	"github.com/joho/godotenv"
)

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
