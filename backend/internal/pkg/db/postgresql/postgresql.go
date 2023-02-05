package database

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var Db *sql.DB

func InitDB() {
	// DB_USER, DB_PASSWORD, DB_NAME := os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE")
	dbinfo := "host=localhost user=postgres password=dbpass dbname=hackernews sslmode=disable"
	db, err := sql.Open("pgx", dbinfo)
	if err != nil {
		log.Panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	Db = db
}

func CloseDB() error {
	return Db.Close()
}

func Migrate() {
	if err := Db.Ping(); err != nil {
		log.Panic(err)
	}
	driver, _ := postgres.WithInstance(Db, &postgres.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/db/migrations/postgresql",
		"postgres",
		driver,
	)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
