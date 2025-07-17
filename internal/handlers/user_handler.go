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
// @Summary      获取用户信息
// @Description  根据用户ID获取用户详细信息
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "用户ID"
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
// @Summary      创建新用户
// @Description  创建一个新的用户记录
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "用户信息"
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

	// 验证用户数据
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
// @Summary      更新用户信息
// @Description  根据用户ID更新用户信息
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int     true  "用户ID"
// @Param        user  body      models.User  true  "用户信息"
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

	// 验证用户数据
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
// @Summary      删除用户
// @Description  根据用户ID删除用户
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "用户ID"
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
// @Summary      搜索用户
// @Description  根据名称搜索用户
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        name    query     string  true  "用户名称搜索模式"
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
// @Summary      列出用户
// @Description  分页列出所有用户
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        limit    query     int  false  "每页数量限制"
// @Param        offset   query     int  false  "偏移量"
// @Success      200  {array}   models.User
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /users/list [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	h.log.Info("Handling ListUsers request")

	// 解析分页参数
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
