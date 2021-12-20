package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"database/sql"
	// use as a sql driver
	_ "github.com/jackc/pgx/v4/stdlib"
)

var db *sql.DB

func connectDB() error {
	var err error
	db, err = sql.Open("pgx", Config.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	return nil
}

// DBPing tests if DB connection is working
func DBPing(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		log.Printf("database is not connected")
		writeStatus(w, http.StatusInternalServerError, "Internal Server Error", false)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		writeStatus(w, http.StatusInternalServerError, "Internal Server Error", false)
		return
	}

	writeStatus(w, http.StatusOK, "OK", false)
}
