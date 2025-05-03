package handlers

import (
	"gopark/internal/db"
	"gopark/internal/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// UserHandler handles user-related requests
type UserHandler struct {
	log *logrus.Logger
	db  *db.DB
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(log *logrus.Logger, db *db.DB) *UserHandler {
	return &UserHandler{log: log, db: db}
}

// GetUser handles GET requests to retrieve user information
func (h *UserHandler) GetUser(c *gin.Context) {
	h.log.Info("Handling GetUser request")
	idParam := c.Query("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.log.Errorf("Invalid ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := h.db.GetUserByID(context.Background(), uint(id))
	if err != nil {
		h.log.Errorf("Failed to retrieve user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser handles POST requests to create a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	h.log.Info("Handling CreateUser request")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.log.Errorf("Invalid request payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.db.CreateUser(context.Background(), &user); err != nil {
		h.log.Errorf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser handles PUT requests to update an existing user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	h.log.Info("Handling UpdateUser request")
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.log.Errorf("Invalid ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.log.Errorf("Invalid request payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user.ID = uint(id)
	if err := h.db.UpdateUser(context.Background(), &user); err != nil {
		h.log.Errorf("Failed to update user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE requests to delete a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	h.log.Info("Handling DeleteUser request")
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.log.Errorf("Invalid ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.db.DeleteUser(context.Background(), uint(id)); err != nil {
		h.log.Errorf("Failed to delete user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
