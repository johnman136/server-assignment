package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// information of mysql
const (
	dbDriver   = "mysql"
	dbUser     = "root"
	dbPassword = "1234"
	dbName     = "testdb"
)

// Send messages to mysql
func SendMessage(w http.ResponseWriter, r *http.Request) {

	// Ensure it is POST request
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Only POST requests are allowed")
		return
	}

	// Retrieve the text data from the request body through query
	sender := r.URL.Query().Get("sender")
	receiver := r.URL.Query().Get("receiver")
	message := r.URL.Query().Get("text")

	// Open a database connection
	db, err := sql.Open(dbDriver, dbUser+":"+dbPassword+"@/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert the text into the database
	cmd := fmt.Sprintf("INSERT INTO data(sender, receiver,message,time_stamp) VALUES (\"%s\",\"%s\",\"%s\",\"%d\")", sender, receiver, message, time.Now().Unix())
	_, err = db.Exec(cmd)
	if err != nil {
		log.Fatal(err)
	}

	// Return a success message
	fmt.Fprintf(w, "Text inserted successfully")
}

// Receive messages from mysql according to "chatroom"
func ReceiveMessage(w http.ResponseWriter, r *http.Request) {

	// Ensure it is GET request
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Only GET requests are allowed")
		return
	}

	// Retrieve the chat data (<user1>:<user2>) from the request body
	chat := strings.Split(r.URL.Query().Get("chat"), ":")
	a, b := chat[0], chat[1]

	// Open a database connection
	db, err := sql.Open(dbDriver, dbUser+":"+dbPassword+"@/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Retrieve the messages from the database
	cmd := fmt.Sprintf("SELECT message FROM data WHERE (receiver = \"%s\" AND sender = \"%s\") OR (sender = \"%s\" AND receiver = \"%s\")", a, b, a, b)
	results, err := db.Query(cmd)
	if err != nil {
		log.Fatal(err)
	}

	counter := 0
	for results.Next() {
		var mesg string

		// for each row, scan the result into our mesg composite object
		err = results.Scan(&mesg)
		if err != nil {
			panic(err.Error())
		}

		counter++
		fmt.Fprintf(w, "Message %d: "+mesg+"\n", counter) // Print out the retrieved messages

	}
}

func main() {
	http.HandleFunc("/api/send", SendMessage)
	http.HandleFunc("/api/pull", ReceiveMessage)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
