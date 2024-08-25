package migrator

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage", "", "path to database")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations folder")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")

	flag.Parse()

	if storagePath == "" {
		panic("storage path is not provided")
	}
	if migrationsPath == "" {
		panic("storage path is not provided")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
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
