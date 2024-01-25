package dbutils

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Karanth1r3/book/pgutil/config"
	"github.com/Karanth1r3/book/pgutil/internal/models"
	_ "github.com/lib/pq"
)

var (
	Hostname = ""
	Port     = "2345"
	Username = ""
	Password = ""
	Database = ""
)

func OpenConnection(cfg config.DB) (*sql.DB, error) {
	// Conn string
	connStrng := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.UserName, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", connStrng)
	if err != nil {
		return nil, fmt.Errorf("could dont open db connection: %w", err)
	}
	return db, nil
}

// If user exists => returns id and nil, else returns -1 (and probably error)
func Exists(username string, conn *sql.DB) (int, error) {
	/*
		// Parsing data for db connection from file
		cfg, err := config.Parse("cfg.yaml")
		if err != nil {
			return -1, fmt.Errorf("Exists() parse: %w", err)
		}
	*/
	username = strings.ToLower(username)
	/*
		// Opening connection to db and deferring it's close to the func end
		db, err := OpenConnection(cfg.DB)
		if err != nil {
			return -1, fmt.Errorf("Exists(): %w", err)
		}
		defer db.Close()
	*/
	// Forming query
	userID := -1
	// Select id column from users table if it contains username
	statement := fmt.Sprintf(`SELECT "id" FROM "users" where 
	username = '%s'`, username)
	rows, err := conn.Query(statement)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return -1, fmt.Errorf("select query failed: %w", err)
		}
		userID = id
	}
	defer rows.Close()
	return userID, nil
}

// Tries to add user to db. Returns id, nil if successfull, -1 and probably err otherwise
func AddUser(data models.UserData, conn *sql.DB) (int, error) {
	/*
		cfg, err := config.Parse("cfg.yaml")
		if err != nil {
			return -1, fmt.Errorf("AddUser() parse error: %w", err)
		}
	*/

	/*
		db, err := OpenConnection(cfg.DB)
		if err != nil {
			return -1, fmt.Errorf("AddUser() open connection error: %w", err)
		}
		defer db.Close()
	*/
	data.Username = strings.ToLower(data.Username)
	// If user already exists => do nothing and return -1,err
	userID, err := Exists(data.Username, conn)
	if userID != -1 {
		return -1, fmt.Errorf("AddUser(): user already exists: %w", err)
	}

	// Insert into users table (FIRST TABLE) concrete username (one value - ($1))
	insertStatement := `insert into "users" ("username") values ($1)`

	_, err = conn.Exec(insertStatement, data.Username)
	if err != nil {
		return -1, fmt.Errorf("AddUser() insert query error (users): %w", err)
	}

	// if user was added to the main table - can proceed; otherwise return -1 and error
	userID, err = Exists(data.Username, conn)
	if userID == -1 {
		return userID, err
	}
	// Inserting into usersdata (SECOND TABLE) new row with 4 values ($1-$4) (using checked userid from added user)
	insertStatement = `insert into "userdata" ("userid", "name", "surname",
	"description") values ($1, $2, $3, $4)`
	_, err = conn.Exec(insertStatement, userID, data.Name, data.Surname, data.Description)
	if err != nil {
		return -1, fmt.Errorf("AddUser() insert query error (userdata): %w", err)
	}
	return userID, nil
}

func DeleteUser(id int, conn *sql.DB) error {
	// Check if user with such id exists
	statement := fmt.Sprintf(`SELECT "username" FROM "users" where id = %d`, id)
	rows, err := conn.Query(statement)
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return fmt.Errorf("DeleteUser() scanning row error: %w", err)
		}
	}
	defer rows.Close()

	if userid, _ := Exists(username, conn); userid != id {
		return fmt.Errorf("user with id %d does not exist", id)
	}
	// Delete statement & execution
	deleteStatement := `delete from "users" where id=$1`
	_, err = conn.Exec(deleteStatement, id)
	if err != nil {
		return fmt.Errorf("DeleteUser() delete query error: %w", err)
	}
	return nil
}

// Gets userdata list from db
func ListUsers(conn *sql.DB) (*[]models.UserData, error) {
	Data := &[]models.UserData{}
	// Get rows from userdata where id from users = id from userdata
	rows, err := conn.Query(`SELECT "id", "username", "name", "surname",
	 "description" FROM "users", "userdata"
	 WHERE users.id = userdata.userid`)
	if err != nil {
		return nil, fmt.Errorf("ListUsers() select query error: %w", err)
	}

	for rows.Next() {
		var id int
		var (
			username, name, surname, description string
		)
		err = rows.Scan(&id, &username, &name, &surname, &description)
		tempData := models.UserData{
			ID:          id,
			Username:    username,
			Name:        name,
			Surname:     surname,
			Description: description,
		}
		*Data = append(*Data, tempData)
		if err != nil {
			return nil, fmt.Errorf("ListUsers() row scan error: %w", err)
		}
	}
	defer rows.Close()
	return Data, nil
}

func UpdateUser(data models.UserData, conn *sql.DB) error {
	// First check if user exists
	userID, err := Exists(data.Username, conn)
	if userID == -1 {
		return fmt.Errorf("UpdateUser(): user does not exist, %w", err)
	}

	// Forming & executing update query with required parameters
	data.ID = userID
	updateStatement := `update "userdata" set "name"=$1, "surname"=$2,
	"description"=$3 where "userid"=$4`
	_, err = conn.Exec(updateStatement, data.Name, data.Surname, data.Description, data.ID)
	if err != nil {
		return fmt.Errorf("UpdateUser() update query execution error: %w", err)
	}
	return nil
}
