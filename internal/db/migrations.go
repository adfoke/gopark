package db

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// MigrationManager handles database migrations
type MigrationManager struct {
	DB  *DB
	Log *logrus.Logger
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *DB, log *logrus.Logger) *MigrationManager {
	return &MigrationManager{
		DB:  db,
		Log: log,
	}
}

// ensureMigrationsTable ensures the migrations table exists
func (m *MigrationManager) ensureMigrationsTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := m.DB.ExecContext(ctx, query)
	if err != nil {
		m.Log.Errorf("Failed to create migrations table: %v", err)
		return err
	}
	return nil
}

// getAppliedMigrations returns the applied migration set
func (m *MigrationManager) getAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	query := `SELECT version FROM schema_migrations;`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		m.Log.Errorf("Failed to query migrations: %v", err)
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			m.Log.Errorf("Failed to scan migration version: %v", err)
			return nil, err
		}
		applied[version] = true
	}

	return applied, nil
}

// recordMigration records an applied migration
func (m *MigrationManager) recordMigration(ctx context.Context, version string) error {
	query := `INSERT INTO schema_migrations (version) VALUES (?);`
	_, err := m.DB.ExecContext(ctx, query, version)
	if err != nil {
		m.Log.Errorf("Failed to record migration %s: %v", version, err)
		return err
	}
	return nil
}

// RunMigrations applies all pending migrations
func (m *MigrationManager) RunMigrations(ctx context.Context, migrationsDir string) error {
	// Ensure the migrations table exists
	if err := m.ensureMigrationsTable(ctx); err != nil {
		return err
	}

	// Fetch applied migrations
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return err
	}

	// Read migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		m.Log.Errorf("Failed to read migrations directory: %v", err)
		return err
	}

	// Filter and sort SQL files
	var migrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrations = append(migrations, file.Name())
		}
	}
	sort.Strings(migrations)

	// Apply migrations
	for _, migration := range migrations {
		// Extract the version (filename prefix)
		version := strings.TrimSuffix(migration, filepath.Ext(migration))

		// Skip already applied migrations
		if applied[version] {
			m.Log.Infof("Migration %s already applied, skipping", version)
			continue
		}

		// Load the migration file
		path := filepath.Join(migrationsDir, migration)
		content, err := os.ReadFile(path)
		if err != nil {
			m.Log.Errorf("Failed to read migration file %s: %v", path, err)
			return err
		}

		// SQLite does not support DDL within transactions; run directly
		// Execute migration
		m.Log.Infof("Applying migration %s", version)
		_, err = m.DB.ExecContext(ctx, string(content))
		if err != nil {
			m.Log.Errorf("Failed to apply migration %s: %v", version, err)
			return err
		}

		// Record the applied migration
		if err := m.recordMigration(ctx, version); err != nil {
			return err
		}

		m.Log.Infof("Successfully applied migration %s", version)
	}

	return nil
}
