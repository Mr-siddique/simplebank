package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/db/util"

	_ "github.com/lib/pq"
)

func main() {
	var err error
	var config util.Config
	config, err = util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	serverAddress := config.ServerAddress
	dbDriver := config.DBDriver
	dbSource := config.DBSource
	conn, err := sql.Open(dbDriver, dbSource) // corrected typo here
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err.Error())
	}
}
