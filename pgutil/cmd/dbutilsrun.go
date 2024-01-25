package main

import (
	"fmt"
	"log"
	"math/rand"

	dbutils "github.com/Karanth1r3/book/pgutil/cmd/DBUtils"
	"github.com/Karanth1r3/book/pgutil/config"
	"github.com/Karanth1r3/book/pgutil/internal/models"
)

var (
	MIN = 0
	MAX = 26
)

// Random return random numbers with set min-max limits
func random(min, max int) int {
	return rand.Intn(max-min) + min
}

// Returns string with randomized characters with specified lenght
func getString(length int64) string {
	var (
		startChar string = "A"
		temp      string = ""
		i         int64  = 1
	)
	for {
		// Getting random character code
		rand := random(MIN, MAX)
		// Selecting random letter (with highcase i guess)
		newChar := string(startChar[0] + byte(rand))
		temp += newChar
		if i == length {
			break
		}
		i++
	}
	return temp
}

func main() {
	// Parsing data with connection parameters from config file
	cfg, err := config.Parse("cfg.yml")
	if err != nil {
		log.Fatal(err)
	}
	// Creating connection
	conn, err := dbutils.OpenConnection(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	// Trying to retreive users list
	data, err := dbutils.ListUsers(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("USERS:")
	for _, elem := range *data {
		fmt.Println(elem)
	}
	randomUsername := getString(5)

	newUser := models.UserData{
		Username:    randomUsername,
		Name:        "Alla",
		Surname:     "Alkina",
		Description: "Who?",
	}
	// Trying  to add user
	id, err := dbutils.AddUser(newUser, conn)
	if id == -1 {
		fmt.Println("error occured while trying to add user: ", err)
	}
	// Trying to delete user
	err = dbutils.DeleteUser(id, conn)
	if err != nil {
		fmt.Println(err)
	}
	// Trying to delete already deleted user out of curiosity
	err = dbutils.DeleteUser(id, conn)
	if err != nil {
		fmt.Println(err)
	}
	//Adding him again
	id, err = dbutils.AddUser(newUser, conn)
	if id == -1 {
		fmt.Println("error occured while trying to add user: ", err)
	}
	//Trying to update user info
	err = dbutils.UpdateUser(newUser, conn)
	if err != nil {
		fmt.Println(err)
	}

}
