package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/TranQuocToan1996/backendMaster/util"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

// https://medium.com/goingogo/why-use-testmain-for-testing-in-go-dafb52b406bc

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(conn)
	testDB = conn
	log.Println("Setup query")

	os.Exit(m.Run())
}
