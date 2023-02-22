package main

import (
	"database/sql"
	"go-exchange/api"
	db "go-exchange/db/sqlc"
	"go-exchange/util"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn := connectToDB(config)

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	runGinServer(config, store)
}

// runGinServer creates and runs a HTTP server with Gin routes
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}

// runDBMigration applies all up migrations
func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance: ", err)
	}

	if err = migration.Up(); err != nil {
		log.Fatal("failed to run migrate up: ", err)
	}

	log.Println("db migrated successfully")
}

// openDB opens and tests db connection
func openDB(dbDriver string, dbSource string) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// connectToDB waits db be ready to connect
func connectToDB(config util.Config) *sql.DB {
	var counts int64

	for {
		connection, err := openDB(config.DBDriver, config.DBSource)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 15 {
			log.Fatal("Can't connect to Postgres: ", err)
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
	}
}
