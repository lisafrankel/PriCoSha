README

**How to run the code**

Download and install go:
* download go here: https://golang.org/dl/
* installation instructions here: https://golang.org/doc/install

The sql driver (used by the backend server) needs to be downloaded/installed:
* in the command line run the following code
..* go get "github.com/go-sql-driver/mysql"

Build executable files, run the following code in the command line:
* go build frontend.go; go build backend.go

Run the servers in the command line:
* ./frontend
* ./backend