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
// @Summary      Get user information
// @Description  Retrieve detailed user information by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      404  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /users [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	h.log.Info("Handling GetUser request")
	idParam := c.Query("id")
	if idParam == "" {
		BadRequest(c, "ID parameter is required", h.log)
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.log.Errorf("Invalid ID format: %v", err)
		BadRequest(c, "Invalid ID format", h.log)
		return
	}

	user, err := h.db.GetUserByID(context.Background(), uint(id))
	if err != nil {
		h.log.Errorf("Failed to retrieve user: %v", err)
		NotFound(c, "User not found", h.log)
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser handles POST requests to create a new user
// @Summary      Create a new user
// @Description  Create a new user record
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User information"
// @Success      201   {object}  models.User
// @Failure      400   {object}  handlers.ErrorResponse
// @Failure      500   {object}  handlers.ErrorResponse
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	h.log.Info("Handling CreateUser request")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.log.Errorf("Invalid request payload: %v", err)
		BadRequest(c, "Invalid request payload", h.log)
		return
	}

	// Validate user data
	if err := user.Validate(); err != nil {
		BadRequest(c, err.Error(), h.log)
		return
	}

	if err := h.db.CreateUser(context.Background(), &user); err != nil {
		h.log.Errorf("Failed to create user: %v", err)
		InternalServerError(c, "Failed to create user", h.log)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser handles PUT requests to update an existing user
// @Summary      Update user information
// @Description  Update user data by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int     true  "User ID"
// @Param        user  body      models.User  true  "User information"
// @Success      200   {object}  models.User
// @Failure      400   {object}  handlers.ErrorResponse
// @Failure      500   {object}  handlers.ErrorResponse
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	h.log.Info("Handling UpdateUser request")
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.log.Errorf("Invalid ID format: %v", err)
		BadRequest(c, "Invalid ID format", h.log)
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.log.Errorf("Invalid request payload: %v", err)
		BadRequest(c, "Invalid request payload", h.log)
		return
	}

	// Validate user data
	if err := user.Validate(); err != nil {
		BadRequest(c, err.Error(), h.log)
		return
	}

	user.ID = uint(id)
	if err := h.db.UpdateUser(context.Background(), &user); err != nil {
		h.log.Errorf("Failed to update user: %v", err)
		InternalServerError(c, "Failed to update user", h.log)
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE requests to delete a user
// @Summary      Delete a user
// @Description  Delete a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	h.log.Info("Handling DeleteUser request")
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.log.Errorf("Invalid ID format: %v", err)
		BadRequest(c, "Invalid ID format", h.log)
		return
	}

	if err := h.db.DeleteUser(context.Background(), uint(id)); err != nil {
		h.log.Errorf("Failed to delete user: %v", err)
		InternalServerError(c, "Failed to delete user", h.log)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// SearchUsers handles GET requests to search users by name
// @Summary      Search users
// @Description  Search for users by name
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        name    query     string  true  "User name search pattern"
// @Success      200  {array}   models.User
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /users/search [get]
func (h *UserHandler) SearchUsers(c *gin.Context) {
	h.log.Info("Handling SearchUsers request")
	namePattern := c.Query("name")
	if namePattern == "" {
		BadRequest(c, "Name search pattern is required", h.log)
		return
	}

	users, err := h.db.SearchUsersByName(context.Background(), namePattern)
	if err != nil {
		h.log.Errorf("Failed to search users: %v", err)
		InternalServerError(c, "Failed to search users", h.log)
		return
	}

	c.JSON(http.StatusOK, users)
}

// ListUsers handles GET requests to list all users with pagination
// @Summary      List users
// @Description  List all users with pagination
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        limit    query     int  false  "Items per page"
// @Param        offset   query     int  false  "Result offset"
// @Success      200  {array}   models.User
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /users/list [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	h.log.Info("Handling ListUsers request")

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		BadRequest(c, "Invalid limit parameter", h.log)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		BadRequest(c, "Invalid offset parameter", h.log)
		return
	}

	users, err := h.db.ListUsers(context.Background(), limit, offset)
	if err != nil {
		h.log.Errorf("Failed to list users: %v", err)
		InternalServerError(c, "Failed to list users", h.log)
		return
	}

	c.JSON(http.StatusOK, users)
}
