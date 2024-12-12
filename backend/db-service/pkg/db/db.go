package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

type DB struct {
	conn *sql.DB
}

func New(connURL string) *DB {
	conn, err := sql.Open("pgx", connURL)
	if err != nil {
		log.Fatal(err)
	}
	err = setupMigrations(conn)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{conn: conn}
}

func (d *DB) Stop() {
	d.conn.Close()
}

func setupMigrations(conn *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	err := goose.Up(conn, "internal/migrations")
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) CreateNewUser(name, about, email string, room, dormNumber int) {
	d.conn.Exec("")
}
