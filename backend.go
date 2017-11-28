package main

import (
	"fmt"
	"net"
	"strings"
	"database/sql"
	_"github.com/go-sql-driver/mysql"

)

func userExists(username string, db *sql.DB) bool {
	fmt.Println("in user exists func")
	fmt.Println("preparing statement")
	stmt, err := db.Prepare("SELECT username FROM Person WHERE username= ?")
	fmt.Println("prepared statement")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()
	rows, _ := stmt.Query(username)
	var count int = 0
	for rows.Next() {
		count++
	}
	if count == 0 {
		return false
	}

	return true

}

// checks to see if password correct for the username
func passwordCorrect(username string, password string, db *sql.DB) bool {
	pw_stmt, _ := db.Prepare("SELECT md5(?)")
	row, _ := pw_stmt.Query(password)
	
	var md5pass string
	for row.Next() {
		err := row.Scan(&md5pass)
		if err != nil {
			fmt.Println(err)
		}
	}


	stmt, err := db.Prepare("SELECT password FROM Person WHERE username= ?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()
	rows, _ := stmt.Query(username)
	var real_password string
	for rows.Next() {
		err := rows.Scan(&real_password)
		if err != nil {
			fmt.Println(err)
		}
	}


	fmt.Println(md5pass, real_password)
	if md5pass == real_password {
		return true
	}

	return false
}

// login handler
func login(username string, password string, db *sql.DB) string {
	fmt.Println("in login func")
	// make sure username is associated with an account & sure password is correct
	if !userExists(username, db) {
		return "601"
	} else if !passwordCorrect(username, password, db) {
		return "602"
	}
	return "200"
}


func createUser(first_name string, last_name string, username string, password string, db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO Person(first_name, last_name, username, password) VALUES (?, ?, ?, md5(?))")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	_, stmt_err := stmt.Query(first_name, last_name, username, password)
	if stmt_err != nil {
		fmt.Println(stmt_err)
	}
}

func makeAccount(first_name string, last_name string, username string, password string, db *sql.DB) string {
	if first_name == "" || last_name == "" || username == "" || password == "" {
		return "606"
	} else {
		if userExists(username, db) {
			return "607"
		} else {
			createUser(first_name, last_name, username, password, db)
		}
	}
	return "200"
}

func handleConn(client net.Conn, db *sql.DB) {
	// read message from client
	message := make([]byte, 1024)
	n, err := client.Read(message)
	if err != nil {
		fmt.Println(err)
	}

	// format message into tokens
	s_message := string(message[:n])
	fmt.Println(s_message)
	message_tokens := strings.Split(s_message, ",")

	if message_tokens[0] == "login" {
		fmt.Println("login request")
		statusCode := login(message_tokens[1], message_tokens[2], db)
		fmt.Fprintf(client, statusCode)
		return
	}
	if message_tokens[0] == "create" {
		fmt.Println("create request")
		statusCode := makeAccount(message_tokens[1], message_tokens[2], message_tokens[3], message_tokens[4], db)
		fmt.Fprintf(client, statusCode)
		return
	}

}


func main() {
	// connect to DB
	fmt.Println("connecting to DB")
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/PriCoSha")
	if err != nil {
		fmt.Println(err)
	}


	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println(err)
	}

	for {
		client, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}

		handleConn(client, db)
	}

}








