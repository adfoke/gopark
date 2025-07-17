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
	// 确保数据库文件所在目录存在
	dbDir := filepath.Dir(cfg.Database.Path)
	if dbDir != "." && dbDir != ".." {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// 连接SQLite数据库
	sqlDB, err := sql.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(25) // SQLite支持的并发连接有限
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
