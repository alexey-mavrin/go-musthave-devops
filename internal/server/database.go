package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn

func connectDB() error {
	var err error
	conn, err = pgx.Connect(context.Background(), Config.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	return nil
}

// DBPing tests if DB connection is working
func DBPing(w http.ResponseWriter, r *http.Request) {
	var resp string
	if conn == nil {
		log.Printf("database is not connected")
		writeStatus(w, http.StatusInternalServerError, "Internal Server Error", false)
		return
	}
	err := conn.QueryRow(context.Background(), "select '1'").Scan(&resp)
	if err != nil {
		writeStatus(w, http.StatusInternalServerError, "Internal Server Error", false)
		return
	}
	writeStatus(w, http.StatusOK, "OK", false)

}
