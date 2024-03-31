package db

import (
	"database/sql"
	"log"
	"os"
	"simplebank/db/util"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	var config util.Config
	config, err = util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	dbDriver := config.DBDriver
	dbSource := config.DBSource
	testDb, err = sql.Open(dbDriver, dbSource) // corrected typo here
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDb)
	os.Exit(m.Run())
}
