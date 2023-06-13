package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct {
	db *sql.DB
}

// information of mysql
const (
	dbDriver   = "mysql"
	dbUser     = "root"
	dbPassword = "1234"
	dbName     = "testdb"
)

// Send chat to mysql with what we get from rpc.SendRequest
func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {

	resp := rpc.NewSendResponse()
	if req.Message == nil {
		//err := fmt.Errorf("Invalid input sanxiao")
		return resp, nil
	}
	Receiver := strings.TrimPrefix(req.Message.GetSender()+":", req.Message.GetChat())

	// Open a database connection
	db, err := sql.Open(dbDriver, dbUser+":"+dbPassword+"@/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert the text into the database
	_, err1 := s.db.Exec("INSERT INTO messages (sender, receiver, text, time) VALUES (?, ?, ?, ?)", req.Message.GetSender(), Receiver, req.Message.GetText(), time.Now().Unix())
	if err1 != nil {
		log.Println("Error storing data in MySQL:", err1)
		resp.Code, resp.Msg = 0, "failed"
		return resp, err1
	}
	resp.Code, resp.Msg = 1, "success"
	return resp, nil

}

// Return responses of chat, nextCursor, etc. according to rpc.PullRequest.cursor, limit, etc. from mysql
func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {

	users := strings.Split(req.Chat, ":")
	Sender, Receiver := users[0], users[1] // Get sender and receiver

	// Open a database connection
	db, err := sql.Open(dbDriver, dbUser+":"+dbPassword+"@/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if req was reversed
	reverse := ""
	if *req.Reverse {
		reverse = "DESC"
	}

	// Retrieve the messages 1 more than LIMIT from the database(e.g. limit = 5 rows, I retrieve 6 rows of data)
	cmd := fmt.Sprintf("SELECT message, sender, receiver, time_stamp FROM data WHERE （(receiver = \"%s\" AND sender = \"%s\") OR (sender = \"%s\" AND receiver = \"%s\")）AND time_stamp >= %d order by id %s LIMIT %d", Sender, Receiver, Sender, Receiver, req.Cursor, reverse, req.Limit+1)
	results, err := db.Query(cmd)
	if err != nil {
		log.Println("Error pulling data from MySQL:", err)
		return nil, err
	}

	// Append data to messages as a slice of data sequentially until hit the limit from where cursor pointed to
	hasMore := false
	var next_cursor int64 = 0
	messages := make([]*rpc.Message, 0)
	var counter int32 = 1
	for results.Next() {
		var message, sender, receiver, time_stamp string

		// For each row, scan the result into composite objects
		err = results.Scan(&message, &sender, &receiver, &time_stamp)
		if err != nil {
			panic(err.Error())
		}

		time_stamp_int, err := strconv.ParseInt(time_stamp, 10, 64)
		if err != nil {
			panic(err.Error())
		}

		// If counter hit the end (LIMIT + 1)
		// - Indicate there are more data (hasMore = true)
		// - Indicate nextCursor
		if counter > req.Limit {
			hasMore = true
			next_cursor = time_stamp_int
			break
		}

		msg := &rpc.Message{
			Chat:     sender + receiver,
			Text:     message,
			Sender:   sender,
			SendTime: time_stamp_int,
		}
		messages = append(messages, msg)
		counter++
	}

	// Return resp
	resp := rpc.NewPullResponse()
	resp.Messages = messages
	resp.Code = 0
	resp.Msg = "success"
	resp.HasMore = &hasMore
	resp.NextCursor = &next_cursor

	return resp, nil

}
