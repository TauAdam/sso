package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	// import drivers
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var dbPath, migrationsPath, migrationsTable string

	flag.StringVar(&dbPath, "database-path", "", "path to database")
	flag.StringVar(&migrationsPath, "migrations-path", "migrations", "path to migrations folder")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")

	flag.Parse()

	if dbPath == "" {
		panic("database path is not provided")
	}
	if migrationsPath == "" {
		panic("migrations folder path is not provided")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", dbPath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no changes")
			return
		}
		panic(err)
	}

	fmt.Println("migrations applied")
}
