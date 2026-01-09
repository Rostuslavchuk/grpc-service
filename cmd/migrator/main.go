package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationOptions struct {
	storagePath     string
	migrationsPath  string
	migrationsTable string
}

func main() {
	migrationOptions := &MigrationOptions{}

	flag.StringVar(&migrationOptions.migrationsPath, "migration-path", "", "get path tp migrations")
	flag.StringVar(&migrationOptions.storagePath, "storage-path", "", "get path to storage path")
	flag.StringVar(&migrationOptions.migrationsTable, "migration-table", "migrations", "name of migrations table")
	flag.Parse()

	if migrationOptions.migrationsPath == "" {
		panic("migration-path not pointed")
	}
	if migrationOptions.storagePath == "" {
		panic("storage-path not pointed")
	}

	m, err := migrate.New(
		"file://"+migrationOptions.migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s",
			migrationOptions.storagePath,
			migrationOptions.migrationsTable,
		),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migration to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migration applied successfully")
}
