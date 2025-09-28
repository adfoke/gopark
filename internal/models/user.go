package models

import (
	"errors"
	"regexp"
	"strings"
)

// User represents the user domain model
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Mail string `json:"mail"`
}

// Validate checks whether the user data is valid
func (u *User) Validate() error {
	// Validate name
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("name is required")
	}

	if len(u.Name) > 255 {
		return errors.New("name is too long (maximum 255 characters)")
	}

	// Validate email
	if strings.TrimSpace(u.Mail) == "" {
		return errors.New("email is required")
	}

	// Perform a basic email format check
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Mail) {
		return errors.New("invalid email format")
	}

	if len(u.Mail) > 255 {
		return errors.New("email is too long (maximum 255 characters)")
	}

	return nil
}
