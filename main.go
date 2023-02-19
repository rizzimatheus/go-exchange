package main

import (
	"database/sql"
	"go-exchange/api"
	db "go-exchange/db/sqlc"
	"go-exchange/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("invalid arguments to connect to db")
	}
	if err := conn.Ping(); err != nil {
		log.Fatal("cannot connect to db")
	}

	log.Printf("conn:%v\nerr:  %v", conn, err)

	store := db.NewStore(conn)

	runGinServer(config, store)
}

// runGinServer creates and runs a HTTP server with Gin routes
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
