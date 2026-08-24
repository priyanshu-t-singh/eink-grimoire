package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	migrations "le-grimoire/db"
	"le-grimoire/internal/constants"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func SetupDatabase() (*sql.DB, error) {
	db, err := Open(constants.SqliteDatabasePath)
	if err != nil {
		return nil, fmt.Errorf("db failed to open: %w", err)
	}
	return db, nil
}

func Open(dbPath string) (*sql.DB, error) {
	// Create Directory if not exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("creating db dir: %w", err)
	}

	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=ON", dbPath)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Ping verifies the connection is actually valid
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	// Configure Goose
	goose.SetBaseFS(migrations.MigrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	slog.Info("running database migrations...")
	if err := goose.Up(db, "migrations"); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	slog.Info("database migrations completed")

	return db, nil
}
