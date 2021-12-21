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

func initDBTable() error {
	_, err := db.Query("CREATE TABLE IF NOT EXISTS gauges (id serial PRIMARY KEY, name VARCHAR (128) UNIQUE NOT NULL, value DOUBLE PRECISION NOT NULL)")
	if err != nil {
		return err
	}
	_, err = db.Query("CREATE TABLE IF NOT EXISTS counters (id serial PRIMARY KEY, name VARCHAR (128) UNIQUE NOT NULL, value BIGINT NOT NULL)")
	return err
}

func storeStatsDB() error {
	_, err := db.Query("DELETE FROM counters")
	if err != nil {
		return err
	}
	_, err = db.Query("DELETE FROM gauges")
	if err != nil {
		return err
	}
	for k, v := range statistics.Gauges {
		_, err = db.Query("INSERT INTO gauges (name, value) VALUES ($1, $2)", k, v)
		log.Print("storing", k)
		if err != nil {
			return err
		}
	}
	for k, v := range statistics.Counters {
		_, err = db.Query("INSERT INTO counters (name, value) VALUES ($1, $2)", k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadStatsDB() error {
	return nil
}

// DBPing tests if DB connection is working
func DBPing(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		log.Printf("database is not connected")
		writeStatus(w, http.StatusInternalServerError, "Internal Server Error", false)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		writeStatus(w, http.StatusInternalServerError, "Internal Server Error", false)
		return
	}

	writeStatus(w, http.StatusOK, "OK", false)
}
