package main

import (
	"fmt"
	"log"
	"net/http"

	Funcs "github.com/advokrat/functions"
	Utils "github.com/advokrat/utilities"

	_ "github.com/mattn/go-sqlite3"

	"github.com/joho/godotenv"
)

//Main Driver Function
func main() {

	http.HandleFunc("/secretpage", Funcs.SecretPageFunc)
	http.HandleFunc("/login", Funcs.LoginFunc)
	http.HandleFunc("/signup", Funcs.SignupFunc)
	http.HandleFunc("/addcoins", Funcs.AddCoinsFunc)
	http.HandleFunc("/transfercoin", Funcs.TransferCoinFunc)
	http.HandleFunc("/getcoins", Funcs.GetCoinsFunc)
	http.HandleFunc("/redeem", Funcs.RedeemCoinsFunc)
	http.HandleFunc("/additems", Funcs.AddItemsFunc)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("An Error Encountered While Loading .env File!!")
	}

	fmt.Println("Listening on Port 8080!!")
	log.Fatal(http.ListenAndServe(":8080", nil))
	defer Utils.Db.Close()
}
