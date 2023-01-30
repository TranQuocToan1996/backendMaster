package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

const (
	// TODO: Move to env var or config file
	dbDriver = "postgres"
	dbSource = "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable"
)

// https://medium.com/goingogo/why-use-testmain-for-testing-in-go-dafb52b406bc

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(conn)
	log.Println("Setup query")

	os.Exit(m.Run())
}

func TestXxx(t *testing.T) {
	log.Println("1")
}
