package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func main() {
	args := os.Args
	if len(args) != 6 {
		fmt.Println("Provide arguments: host port username password db")
		return
	}
	// Connection string for db
	host := args[1]
	p := args[2]
	user := args[3]
	pass := args[4]
	database := args[5]

	port, err := strconv.Atoi(p)
	if err != nil {
		fmt.Println("not a valid port number:", err)
		return
	}

	conn := fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s
	sslmode=disable`, host, port, user, pass, database)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("sql.Open():", err)
		return
	}
	fmt.Println(db)
	defer db.Close()
	// Retreive all dbs
	rows, err := db.Query(`SELECT "datname" FROM pg_database
WHERE datistemplate = false`)
	if err != nil {
		fmt.Println("Query1:", err)
		return
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Scan", err)
			return
		}
		fmt.Println("*", name)
	}
	defer rows.Close()

	// Get all tables from current db
	query := `SELECT table_name FROM information_schema.tables WHERE
	table_schema = 'public' ORDER BY table_name`
	rows, err = db.Query(query)
	if err != nil {
		fmt.Println("Query2", err)
		return
	}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Scan", err)
			return
		}
		fmt.Println("+T", name)
	}
	defer rows.Close()
}
