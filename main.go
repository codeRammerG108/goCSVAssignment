package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	DBconnection "github.com/codeRammerG108/goCSVAssignment/db"
)

var DB *sql.DB

func main() {
	var err error
	fmt.Println("Go-Lang Assignment on CSV")
	DB, err := DBconnection.DBinit()
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("\nConnection Established")
	}
	defer DB.Close()

}
