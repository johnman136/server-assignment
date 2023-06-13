package main

import (
	"log"

	"database/sql"

	rpc "github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc/imservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Open a database connection
	db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/testdb")
	if err != nil {
		log.Fatal("Error connecting to MySQL:", err)
	}
	defer db.Close()

	svr := rpc.NewServer(&IMServiceImpl{db: db}, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "demo.rpc.server",
	}))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}

}
