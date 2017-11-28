README

**How to run the code using command line**

Download and install go:
* download go here: https://golang.org/dl/
* installation instructions here: https://golang.org/doc/install

The sql driver (used by the backend server) needs to be downloaded/installed:
* go get "github.com/go-sql-driver/mysql"

Build executable files, run the following code:
* go build frontend.go; go build backend.go

Run the servers:
* ./frontend
* ./backend

**Use the webiste**

First, make sure the servers are running.
In order to login, go to http://localhost:8080/login
In order to create an account, go to http://localhost:8080/create

Logging in will direct the user to http://localhost:8080/main - where a user has the option logout.