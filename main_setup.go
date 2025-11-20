package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/adamjames870/chirpy/internal/database"
)

func (state *apiState) LoadState() error {
	state.dbQueries = database.New(loadDb())
	state.platform = loadPlatform()
	return nil
}

func loadDb() *sql.DB {
	dbUrl := os.Getenv("DB_URL")
	db, errDb := sql.Open("postgres", dbUrl)
	if errDb != nil {
		fmt.Println("Unable to load DB: " + errDb.Error())
		os.Exit(1)
	}
	return db
}

func loadPlatform() string {
	return os.Getenv("PLATFORM")
}
