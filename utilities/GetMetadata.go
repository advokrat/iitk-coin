package utilities

import (

	//Func "github.com/advokrat/functions"
	//Utils "iitk-coin/utilities"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
)

//Returns ROll NO. of the User associated with the Token
func GetMetadata(user_token string) (string, string, error) {
	token, err := AuthJWT(user_token)
	if err != nil {
		return " ", " ", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		roll_no, _ := claims["user_roll_no"].(string)
		account_type, _ := claims["accountType"].(string)
		return roll_no, account_type, err
	}

	return " ", " ", err

}
