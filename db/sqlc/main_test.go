package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/rs/zerolog/log"

	"github.com/TranQuocToan1996/backendMaster/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

// https://medium.com/goingogo/why-use-testmain-for-testing-in-go-dafb52b406bc

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	testQueries = New(conn)
	testDB = conn
	log.Info().Msg("Setup query")

	os.Exit(m.Run())
}
