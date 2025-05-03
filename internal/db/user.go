package db

import (
	"gopark/internal/models"
	"golang.org/x/net/context"
)

// CreateUser inserts a new user into the database
func (db *DB) CreateUser(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (name, mail) VALUES ($1, $2) RETURNING id"
	err := db.Pool.QueryRow(ctx, query, user.Name, user.Mail).Scan(&user.ID)
	if err != nil {
		db.Log.Errorf("Failed to create user: %v", err)
		return err
	}
	db.Log.Infof("Created user with ID %d", user.ID)
	return nil
}

// GetUserByID retrieves a user by their ID
func (db *DB) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	query := "SELECT id, name, mail FROM users WHERE id = $1"
	user := &models.User{}
	err := db.Pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Mail)
	if err != nil {
		db.Log.Errorf("Failed to get user by ID %d: %v", id, err)
		return nil, err
	}
	return user, nil
}

// UpdateUser updates an existing user in the database
func (db *DB) UpdateUser(ctx context.Context, user *models.User) error {
	query := "UPDATE users SET name = $2, mail = $3 WHERE id = $1"
	_, err := db.Pool.Exec(ctx, query, user.ID, user.Name, user.Mail)
	if err != nil {
		db.Log.Errorf("Failed to update user ID %d: %v", user.ID, err)
		return err
	}
	db.Log.Infof("Updated user with ID %d", user.ID)
	return nil
}

// DeleteUser deletes a user from the database by ID
func (db *DB) DeleteUser(ctx context.Context, id uint) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		db.Log.Errorf("Failed to delete user ID %d: %v", id, err)
		return err
	}
	db.Log.Infof("Deleted user with ID %d", id)
	return nil
}
