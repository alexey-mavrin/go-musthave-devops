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

func runSQLStatement(s string) error {
	rows, err := db.Query(s)
	if err != nil {
		return err
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func initDBTable() error {
	err := runSQLStatement("CREATE TABLE IF NOT EXISTS gauges (id serial PRIMARY KEY, name VARCHAR (128) UNIQUE NOT NULL, value DOUBLE PRECISION NOT NULL)")
	if err != nil {
		return err
	}

	err = runSQLStatement("CREATE TABLE IF NOT EXISTS counters (id serial PRIMARY KEY, name VARCHAR (128) UNIQUE NOT NULL, value BIGINT NOT NULL)")
	if err != nil {
		return err
	}

	return nil
}

func storeStatsDB() error {
	err := runSQLStatement("DELETE FROM counters")
	if err != nil {
		return err
	}
	err = runSQLStatement("DELETE FROM gauges")
	if err != nil {
		return err
	}
	for k, v := range statistics.Gauges {
		err = storeGaugeDB(k, v)
		if err != nil {
			return err
		}
	}
	for k, v := range statistics.Counters {
		err = storeCounterDB(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func storeGaugeDB(name string, gauge float64) error {
	rows, err := db.Query("INSERT INTO gauges (name, value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE set value = $2", name, gauge)
	if err != nil {
		return err
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return err
	}
	// log.Printf("gauge %s = %v stored in db", name, gauge)
	return nil
}

func storeCounterDB(name string, counter int64) error {
	rows, err := db.Query("INSERT INTO counters (name, value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET value = $2", name, counter)
	if err != nil {
		return err
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return err
	}
	// log.Printf("counter %s = %v stored in db", name, counter)
	return nil
}

func loadStatsDB() error {
	var name string
	var gauge float64
	var counter int64

	mu.Lock()
	defer mu.Unlock()

	gRows, err := db.Query("SELECT name, value FROM gauges")
	if err != nil {
		return err
	}
	defer gRows.Close()
	for gRows.Next() {
		if err = gRows.Scan(&name, &gauge); err != nil {
			log.Print(err)
			return err
		}
		statistics.Gauges[name] = gauge
	}
	if err = gRows.Err(); err != nil {
		return err
	}

	cRows, err := db.Query("SELECT name, value FROM counters")
	if err != nil {
		return err
	}
	defer cRows.Close()
	for cRows.Next() {
		if err = cRows.Scan(&name, &counter); err != nil {
			log.Print(err)
			return err
		}
		statistics.Counters[name] = counter
	}
	if err = cRows.Err(); err != nil {
		return err
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
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		writeStatus(w, http.StatusInternalServerError, "Internal Server Error", false)
		return
	}

	writeStatus(w, http.StatusOK, "OK", false)
}
