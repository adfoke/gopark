package db

import (
	"fmt"
	"gopark/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// DB holds the database connection pool
type DB struct {
	Pool *pgxpool.Pool
	Log  *logrus.Logger
}

// NewDB initializes a new database connection pool
func NewDB(cfg config.Config, log *logrus.Logger) (*DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	poolConfig, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	db := &DB{
		Pool: poolConfig,
		Log:  log,
	}

	log.Infof("Database connection established to %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	return db, nil
}

// Close closes the database connection pool
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		db.Log.Info("Database connection pool closed")
	}
}
