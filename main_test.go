package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var db *sql.DB

func TestMain(m *testing.M) {
	initialize()
	code := m.Run()
	clearDB()
	os.Exit(code)
}

func initialize() {
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", "5433", "zanzibert", "nekineki", "test")

	db, _ = sql.Open("postgres", connection)
}

func clearDB() {
	db.Exec("delete from users")
	db.Exec("alter sequence users_id_seq restart with 1")
	db.Exec("delete from groups")
	db.Exec("alter sequence groups_id_seq restart with 1")
}
