package db

import (
	"database/sql"
	"gopark/internal/models"

	"golang.org/x/net/context"
)

// CreateUser inserts a new user into the database
func (db *DB) CreateUser(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (name, mail) VALUES (?, ?)"
	result, err := db.ExecContext(ctx, query, user.Name, user.Mail)
	if err != nil {
		db.Log.Errorf("Failed to create user: %v", err)
		return err
	}

	// Retrieve auto-incremented ID
	id, err := result.LastInsertId()
	if err != nil {
		db.Log.Errorf("Failed to get last insert ID: %v", err)
		return err
	}

	user.ID = uint(id)
	db.Log.Infof("Created user with ID %d", user.ID)
	return nil
}

// GetUserByID retrieves a user by their ID
func (db *DB) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	query := "SELECT id, name, mail FROM users WHERE id = ?"
	user := &models.User{}
	err := db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Mail)
	if err != nil {
		if err == sql.ErrNoRows {
			db.Log.Infof("No user found with ID %d", id)
		} else {
			db.Log.Errorf("Failed to get user by ID %d: %v", id, err)
		}
		return nil, err
	}
	return user, nil
}

// UpdateUser updates an existing user in the database
func (db *DB) UpdateUser(ctx context.Context, user *models.User) error {
	query := "UPDATE users SET name = ?, mail = ? WHERE id = ?"
	result, err := db.ExecContext(ctx, query, user.Name, user.Mail, user.ID)
	if err != nil {
		db.Log.Errorf("Failed to update user ID %d: %v", user.ID, err)
		return err
	}

	// Verify that a row was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.Log.Errorf("Failed to get rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		db.Log.Warnf("No user found with ID %d for update", user.ID)
		return sql.ErrNoRows
	}

	db.Log.Infof("Updated user with ID %d", user.ID)
	return nil
}

// DeleteUser deletes a user from the database by ID
func (db *DB) DeleteUser(ctx context.Context, id uint) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		db.Log.Errorf("Failed to delete user ID %d: %v", id, err)
		return err
	}

	// Verify that a row was deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.Log.Errorf("Failed to get rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		db.Log.Warnf("No user found with ID %d for deletion", id)
		return sql.ErrNoRows
	}

	db.Log.Infof("Deleted user with ID %d", id)
	return nil
}

// SearchUsersByName searches for users by name pattern
func (db *DB) SearchUsersByName(ctx context.Context, namePattern string) ([]*models.User, error) {
	// SQLite uses LIKE instead of ILIKE; apply COLLATE NOCASE for case-insensitive matching
	query := "SELECT id, name, mail FROM users WHERE name LIKE ? COLLATE NOCASE ORDER BY id LIMIT 100"
	rows, err := db.QueryContext(ctx, query, "%"+namePattern+"%")
	if err != nil {
		db.Log.Errorf("Failed to search users by name pattern '%s': %v", namePattern, err)
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Mail); err != nil {
			db.Log.Errorf("Failed to scan user row: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		db.Log.Errorf("Error iterating user rows: %v", err)
		return nil, err
	}

	db.Log.Infof("Found %d users matching name pattern '%s'", len(users), namePattern)
	return users, nil
}

// ListUsers retrieves all users with pagination
func (db *DB) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	if limit <= 0 {
		limit = 10 // Default limit
}
	if limit > 100 {
		limit = 100 // Maximum limit
	}
	if offset < 0 {
		offset = 0
	}

	query := "SELECT id, name, mail FROM users ORDER BY id LIMIT ? OFFSET ?"
	rows, err := db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		db.Log.Errorf("Failed to list users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Mail); err != nil {
			db.Log.Errorf("Failed to scan user row: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		db.Log.Errorf("Error iterating user rows: %v", err)
		return nil, err
	}

	db.Log.Infof("Listed %d users (limit: %d, offset: %d)", len(users), limit, offset)
	return users, nil
}
