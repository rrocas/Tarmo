package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	// Ensure the parent directory exists, since SQLite won't create it
	if dir := filepath.Dir(dbPath); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create db directory %q: %w", dir, err)
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database file %q: %w", dbPath, err)
	}

	// sql.Open doesn't actually connect — verify the file/connection works
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable WAL mode for better concurrency (allows concurrent reads and writes)
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}
	// Set busy timeout to 5 seconds to handle lock contention
	if _, err := db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		return nil, fmt.Errorf("failed to set busy timeout: %w", err)
	}
	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	dbInstance := &SQLiteDB{db: db}
	if err := dbInstance.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}
	return dbInstance, nil
}

func (r *SQLiteDB) init() error {
	query := `
  CREATE TABLE IF NOT EXISTS templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    quantity REAL,
    unit TEXT,
    difficulty INTEGER
  );

  CREATE TABLE IF NOT EXISTS template_resources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL,
    resource_id INTEGER NOT NULL,
    quantity REAL,
    unit TEXT,
    FOREIGN KEY (template_id) REFERENCES templates(id) ON DELETE CASCADE,
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE
  );

  CREATE TABLE IF NOT EXISTS steps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL,
    step_order INTEGER NOT NULL,
    name TEXT NOT NULL,
    instructions TEXT,
    FOREIGN KEY (template_id) REFERENCES templates(id) ON DELETE CASCADE
  );

  CREATE TABLE IF NOT EXISTS resources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    price INTEGER,
    base_quantity REAL,
    base_unit TEXT
  );

  
  `
	_, err := r.db.Exec(query)
	return err
}
