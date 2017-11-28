package main

import (
	"fmt"
	"net/http"
	"net"
	"time"
	"html/template"

)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login handler")
	userCookie, _ := r.Cookie("SessionID")
	if userCookie != nil {
		http.Redirect(w, r, "/home", 302)
		return
	}
	if r.Method == "GET" {
		fmt.Println("login get req")
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		fmt.Println("login post req")
		r.ParseForm()
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")	


		conn, _ := net.Dial("tcp", "localhost:8081")
		defer conn.Close()
		client_msg := "login," + username +"," + password
		fmt.Fprintf(conn, client_msg)
		
		response := make([]byte, 1024)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println(err)
			}
		s_message := string(response[:n])

		if s_message == "601" {
			fmt.Fprintf(w, "User does not exist")
		} else if s_message == "602" {
			fmt.Fprintf(w, "Incorrect password")
		} else {
			cookieValue := r.PostFormValue("username")
			expire := time.Now().Add(1 * time.Hour)
			userCookie := http.Cookie{Name: "SessionID", Value: cookieValue, Expires: expire}
			http.SetCookie(w, &userCookie)
			http.Redirect(w, r, "/login", 302)
		}
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("create handler")
	userCookie, _ := r.Cookie("SessionID")
	if userCookie != nil {
		http.Redirect(w, r, "/home", 302)
		return
	}
	if r.Method == "GET" {
		fmt.Println("create get request")
		t, _ := template.ParseFiles("create.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		fmt.Println("create post request")
		r.ParseForm()
		
		// connect to server & define variables that will be sent to server	
		conn, _ := net.Dial("tcp", "localhost:8081")
		defer conn.Close()
		first_name := r.PostFormValue("first_name")
		last_name := r.PostFormValue("last_name")
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		
		// compose message to client, of request, username, name, email & password
		client_msg := "create," + first_name + "," + last_name + "," + username + "," + password
		
		// send message to server
		fmt.Fprintf(conn, client_msg)
		//_, err := bufio.NewReader(conn).ReadString('\n')
		response := make([]byte, 1024)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println(err)
		}

		s_message := string(response[:n])
		if s_message == "200" {
			http.Redirect(w, r, "/login", 302)
			return
		}

		if s_message == "607" {
			fmt.Fprintf(w, "Username already exists")
		} else if s_message == "606" {
			fmt.Fprintf(w, "Invalid entry")
		} 
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	userCookie, _ := r.Cookie("SessionID")
	if userCookie == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if r.Method == "GET" {
		t, _ := template.ParseFiles("main.html")
		t.Execute(w, nil)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	userCookie, _ := r.Cookie("SessionID")
	if userCookie == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if r.Method == "GET" {
		t, _ := template.ParseFiles("logout.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {

		// expire cookie & redirect them back to login page
		username := userCookie.Value
		cookieValue := username
		expire := time.Now()
		userCookie := http.Cookie{Name: "SessionID", Value: cookieValue, Expires: expire}
		http.SetCookie(w, &userCookie)
		http.Redirect(w, r, "/login", 302)
	}
}


func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/create", create)
	http.HandleFunc("/home", home)
	http.HandleFunc("/logout", logout)
	http.ListenAndServe(":8080", nil)
}