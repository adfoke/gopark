package db

import (
	"database/sql"
	"fmt"
	"gopark/config"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// DB holds the database connection
type DB struct {
	DB  *sql.DB
	Log *logrus.Logger
}

// NewDB initializes a new database connection
func NewDB(cfg config.Config, log *logrus.Logger) (*DB, error) {
	// Ensure the database directory exists
	dbDir := filepath.Dir(cfg.Database.Path)
	if dbDir != "." && dbDir != ".." {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Open the SQLite database
	sqlDB, err := sql.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	// Verify the connection
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Configure the connection pool
	sqlDB.SetMaxOpenConns(25) // SQLite allows limited concurrent connections
	sqlDB.SetMaxIdleConns(5)

	db := &DB{
		DB:  sqlDB,
		Log: log,
	}

	log.Infof("Database connection established to %s", cfg.Database.Path)
	return db, nil
}

// Close closes the database connection
func (db *DB) Close() {
	if db.DB != nil {
		db.DB.Close()
		db.Log.Info("Database connection closed")
	}
}

// ExecContext executes a query without returning any rows
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.DB.ExecContext(ctx, query, args...)
}

// QueryContext executes a query that returns rows
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.DB.QueryContext(ctx, query, args...)
}

// QueryRowContext executes a query that returns a single row
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRowContext(ctx, query, args...)
}
