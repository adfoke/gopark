package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"gopark/internal/db"
	"gopark/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

// 创建一个模拟的数据库
type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	// 模拟ID赋值
	user.ID = 1
	return args.Error(0)
}

func (m *MockDB) UpdateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockDB) DeleteUser(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// 测试辅助函数
func setupTest() (*gin.Engine, *MockDB, *logrus.Logger) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockDB := new(MockDB)
	log := logrus.New()
	log.SetOutput(bytes.NewBuffer(nil)) // 禁止日志输出
	return r, mockDB, log
}

// 测试 GetUser
func TestGetUser(t *testing.T) {
	r, mockDB, log := setupTest()
	handler := &UserHandler{log: log, db: &db.DB{}}

	// 替换handler中的db为mockDB
	handler.db = &db.DB{} // 这里只是为了类型匹配，实际操作会使用mockDB

	// 注册路由
	r.GET("/user", handler.GetUser)

	// 测试用例1: 成功获取用户
	t.Run("Success", func(t *testing.T) {
		mockUser := &models.User{ID: 1, Name: "Test User", Mail: "test@example.com"}
		mockDB.On("GetUserByID", mock.Anything, uint(1)).Return(mockUser, nil).Once()

		// 创建请求
		req, _ := http.NewRequest("GET", "/user?id=1", nil)
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mockUser.ID, response.ID)
		assert.Equal(t, mockUser.Name, response.Name)
		assert.Equal(t, mockUser.Mail, response.Mail)

		mockDB.AssertExpectations(t)
	})

	// 测试用例2: 缺少ID参数
	t.Run("Missing ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/user", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID parameter is required", response.Message)
	})

	// 测试用例3: 无效的ID格式
	t.Run("Invalid ID Format", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/user?id=invalid", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid ID format", response.Message)
	})

	// 测试用例4: 用户不存在
	t.Run("User Not Found", func(t *testing.T) {
		mockDB.On("GetUserByID", mock.Anything, uint(999)).Return(nil, errors.New("user not found")).Once()

		req, _ := http.NewRequest("GET", "/user?id=999", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User not found", response.Message)

		mockDB.AssertExpectations(t)
	})
}

// 测试 CreateUser
func TestCreateUser(t *testing.T) {
	r, mockDB, log := setupTest()
	handler := &UserHandler{log: log, db: &db.DB{}}

	// 注册路由
	r.POST("/user", handler.CreateUser)

	// 测试用例1: 成功创建用户
	t.Run("Success", func(t *testing.T) {
		mockUser := &models.User{Name: "New User", Mail: "new@example.com"}
		mockDB.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()

		// 创建请求
		userJSON, _ := json.Marshal(mockUser)
		req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 验证响应
		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID) // 模拟函数设置ID为1
		assert.Equal(t, mockUser.Name, response.Name)
		assert.Equal(t, mockUser.Mail, response.Mail)

		mockDB.AssertExpectations(t)
	})

	// 测试用例2: 无效的请求体
	t.Run("Invalid Request Body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request payload", response.Message)
	})

	// 测试用例3: 缺少必填字段
	t.Run("Missing Required Fields", func(t *testing.T) {
		mockUser := &models.User{Name: "", Mail: ""}

		userJSON, _ := json.Marshal(mockUser)
		req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Name is required", response.Message)
	})

	// 测试用例4: 数据库错误
	t.Run("Database Error", func(t *testing.T) {
		mockUser := &models.User{Name: "Error User", Mail: "error@example.com"}
		mockDB.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("database error")).Once()

		userJSON, _ := json.Marshal(mockUser)
		req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to create user", response.Message)

		mockDB.AssertExpectations(t)
	})
}

// 测试 UpdateUser
func TestUpdateUser(t *testing.T) {
	r, mockDB, log := setupTest()
	handler := &UserHandler{log: log, db: &db.DB{}}

	// 注册路由
	r.PUT("/user/:id", handler.UpdateUser)

	// 测试用例1: 成功更新用户
	t.Run("Success", func(t *testing.T) {
		mockUser := &models.User{ID: 1, Name: "Updated User", Mail: "updated@example.com"}
		mockDB.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()

		// 创建请求
		userJSON, _ := json.Marshal(mockUser)
		req, _ := http.NewRequest("PUT", "/user/1", bytes.NewBuffer(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, mockUser.Name, response.Name)
		assert.Equal(t, mockUser.Mail, response.Mail)

		mockDB.AssertExpectations(t)
	})

	// 测试用例2: 无效的ID格式
	t.Run("Invalid ID Format", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/user/invalid", bytes.NewBuffer([]byte(`{"name":"Test","mail":"test@example.com"}`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid ID format", response.Message)
	})

	// 测试用例3: 无效的请求体
	t.Run("Invalid Request Body", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/user/1", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request payload", response.Message)
	})

	// 测试用例4: 缺少必填字段
	t.Run("Missing Required Fields", func(t *testing.T) {
		mockUser := &models.User{ID: 1, Name: "", Mail: ""}

		userJSON, _ := json.Marshal(mockUser)
		req, _ := http.NewRequest("PUT", "/user/1", bytes.NewBuffer(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Name is required", response.Message)
	})

	// 测试用例5: 数据库错误
	t.Run("Database Error", func(t *testing.T) {
		mockUser := &models.User{ID: 1, Name: "Error User", Mail: "error@example.com"}
		mockDB.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("database error")).Once()

		userJSON, _ := json.Marshal(mockUser)
		req, _ := http.NewRequest("PUT", "/user/1", bytes.NewBuffer(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to update user", response.Message)

		mockDB.AssertExpectations(t)
	})
}

// 测试 DeleteUser
func TestDeleteUser(t *testing.T) {
	r, mockDB, log := setupTest()
	handler := &UserHandler{log: log, db: &db.DB{}}

	// 注册路由
	r.DELETE("/user/:id", handler.DeleteUser)

	// 测试用例1: 成功删除用户
	t.Run("Success", func(t *testing.T) {
		mockDB.On("DeleteUser", mock.Anything, uint(1)).Return(nil).Once()

		// 创建请求
		req, _ := http.NewRequest("DELETE", "/user/1", nil)
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User deleted successfully", response["message"])

		mockDB.AssertExpectations(t)
	})

	// 测试用例2: 无效的ID格式
	t.Run("Invalid ID Format", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/user/invalid", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid ID format", response.Message)
	})

	// 测试用例3: 数据库错误
	t.Run("Database Error", func(t *testing.T) {
		mockDB.On("DeleteUser", mock.Anything, uint(999)).Return(errors.New("database error")).Once()

		req, _ := http.NewRequest("DELETE", "/user/999", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to delete user", response.Message)

		mockDB.AssertExpectations(t)
	})
}
