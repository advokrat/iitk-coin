package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, _ :=
		sql.Open("sqlite3", "./students.db")
	createTable, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS students (rollno INTEGER PRIMARY KEY, name TEXT)")
	createTable.Exec()
	insertEntry, _ :=
		database.Prepare("INSERT INTO students (rollno, name) VALUES (?, ?)")
	insertEntry.Exec(190436, "Kratik Agrawal")
	rows, _ :=
		database.Query("SELECT rollno, name FROM students")
	var rollno int
	var name string
	for rows.Next() {
		rows.Scan(&rollno, &name)
		fmt.Println(strconv.Itoa(rollno) + ": " + name)
	}
}
