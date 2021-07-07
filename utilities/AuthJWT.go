package utilities

import (
	"fmt"
	"log"
	"os"

	//Func "github.com/advokrat/functions"
	//Utils "iitk-coin/utilities"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"

	"github.com/joho/godotenv"
)

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
