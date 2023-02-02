package main

import (
	"database/sql"
	"log"

	"github.com/TranQuocToan1996/backendMaster/api"
	db "github.com/TranQuocToan1996/backendMaster/db/sqlc"
	"github.com/TranQuocToan1996/backendMaster/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(config)

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
