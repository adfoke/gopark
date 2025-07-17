package db

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// MigrationManager 处理数据库迁移
type MigrationManager struct {
	DB  *DB
	Log *logrus.Logger
}

// NewMigrationManager 创建一个新的迁移管理器
func NewMigrationManager(db *DB, log *logrus.Logger) *MigrationManager {
	return &MigrationManager{
		DB:  db,
		Log: log,
	}
}

// ensureMigrationsTable 确保迁移表存在
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

// getAppliedMigrations 获取已应用的迁移列表
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

// recordMigration 记录已应用的迁移
func (m *MigrationManager) recordMigration(ctx context.Context, version string) error {
	query := `INSERT INTO schema_migrations (version) VALUES (?);`
	_, err := m.DB.ExecContext(ctx, query, version)
	if err != nil {
		m.Log.Errorf("Failed to record migration %s: %v", version, err)
		return err
	}
	return nil
}

// RunMigrations 运行所有未应用的迁移
func (m *MigrationManager) RunMigrations(ctx context.Context, migrationsDir string) error {
	// 确保迁移表存在
	if err := m.ensureMigrationsTable(ctx); err != nil {
		return err
	}

	// 获取已应用的迁移
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return err
	}

	// 读取迁移文件
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		m.Log.Errorf("Failed to read migrations directory: %v", err)
		return err
	}

	// 过滤并排序SQL文件
	var migrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrations = append(migrations, file.Name())
		}
	}
	sort.Strings(migrations)

	// 应用迁移
	for _, migration := range migrations {
		// 提取版本号（文件名前缀）
		version := strings.TrimSuffix(migration, filepath.Ext(migration))

		// 检查是否已应用
		if applied[version] {
			m.Log.Infof("Migration %s already applied, skipping", version)
			continue
		}

		// 读取迁移文件
		path := filepath.Join(migrationsDir, migration)
		content, err := os.ReadFile(path)
		if err != nil {
			m.Log.Errorf("Failed to read migration file %s: %v", path, err)
			return err
		}

		// SQLite不支持事务中的DDL语句，所以我们不使用事务
		// 执行迁移
		m.Log.Infof("Applying migration %s", version)
		_, err = m.DB.ExecContext(ctx, string(content))
		if err != nil {
			m.Log.Errorf("Failed to apply migration %s: %v", version, err)
			return err
		}

		// 记录迁移
		if err := m.recordMigration(ctx, version); err != nil {
			return err
		}

		m.Log.Infof("Successfully applied migration %s", version)
	}

	return nil
}
