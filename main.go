package main

import (
	"database/sql"
	"log"

	"github.com/TranQuocToan1996/backendMaster/api"
	db "github.com/TranQuocToan1996/backendMaster/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	// TODO: Move to env var or config file
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	server.Start(serverAddress)
}
