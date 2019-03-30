package pgstorage_test

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq" // init pg driver
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
)

const existingDBName = "coinsph"

var mutex = &sync.Mutex{}

// prepareDB creates a temporary database with unique name
func prepareDB(t *testing.T) (*sql.DB, func()) {
	dbName, err := createTempDB(existingDBName)
	if err != nil {
		t.Fatalf("unable to create temp db: %s", err)
	}

	db, err := sql.Open("postgres", connectionString(dbName))
	if err != nil {
		t.Fatalf("unable to connect to temp db: %s", err)
	}

	if err = applyDBMigrations(db); err != nil {
		defer db.Close()
		t.Fatalf("unable to apply migrations: %s", err)
	}

	return db, func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close txdb connection: %s", err)
		}
		if err := dropTestDB(existingDBName, dbName); err != nil {
			t.Fatalf("unable to drop temp db %s: %s", dbName, err)
		}
	}
}

func createTempDB(existingDBName string) (newDBName string, err error) {
	newDBName = fmt.Sprintf("coinsph_test_%d", time.Now().UnixNano())

	db, err := sql.Open("postgres", connectionString(existingDBName))
	if err != nil {
		return newDBName, err
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil && err == nil {
			err = errors.Wrap(closeErr, "unable to close main DB")
		}
	}()

	_, err = db.Exec(`CREATE DATABASE "` + newDBName + `"`)
	return newDBName, err
}

func dropTestDB(existingDBName, testDBName string) error {
	db, err := sql.Open("postgres", connectionString(existingDBName))
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil && err == nil {
			err = errors.Wrap(closeErr, "unable to close main DB")
		}
	}()

	_, err = db.Exec(`DROP DATABASE "` + testDBName + `"`)
	return err
}

func applyDBMigrations(db *sql.DB) error {
	mutex.Lock()
	migrate.SetTable("migrations")
	migrations := &migrate.FileMigrationSource{
		Dir: "../../migrations",
	}
	_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	mutex.Unlock()
	return err
}

func connectionString(dbName string) string {
	return fmt.Sprintf("postgres://localhost/%s?sslmode=disable", dbName)
}
