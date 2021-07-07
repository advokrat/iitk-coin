package utilities

import (
	"database/sql"
	"strconv"

	//Func "github.com/advokrat/functions"
	//Utils "iitk-coin/utilities"

	_ "github.com/mattn/go-sqlite3"
)

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
