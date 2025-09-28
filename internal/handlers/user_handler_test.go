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

// MockDB provides a mocked database implementation
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
	// Simulate ID assignment
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

// setupTest prepares the gin router, mock DB, and logger
func setupTest() (*gin.Engine, *MockDB, *logrus.Logger) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockDB := new(MockDB)
	log := logrus.New()
	log.SetOutput(bytes.NewBuffer(nil)) // Disable logging output
	return r, mockDB, log
}

// TestGetUser exercises the GetUser handler
func TestGetUser(t *testing.T) {
	r, mockDB, log := setupTest()
	handler := &UserHandler{log: log, db: &db.DB{}}

// Replace handler DB with mockDB
handler.db = &db.DB{} // Placeholder type; mockDB handles interactions

	// Register route
	r.GET("/user", handler.GetUser)

	// Test case 1: successfully retrieve a user
	t.Run("Success", func(t *testing.T) {
		mockUser := &models.User{ID: 1, Name: "Test User", Mail: "test@example.com"}
		mockDB.On("GetUserByID", mock.Anything, uint(1)).Return(mockUser, nil).Once()

		// Create request
		req, _ := http.NewRequest("GET", "/user?id=1", nil)
		w := httptest.NewRecorder()

		// Execute request
		r.ServeHTTP(w, req)

		// Validate response
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mockUser.ID, response.ID)
		assert.Equal(t, mockUser.Name, response.Name)
		assert.Equal(t, mockUser.Mail, response.Mail)

		mockDB.AssertExpectations(t)
	})

	// Test case 2: missing ID parameter
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

	// Test case 3: invalid ID format
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

	// Test case 4: user not found
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

// TestCreateUser exercises the CreateUser handler
func TestCreateUser(t *testing.T) {
	r, mockDB, log := setupTest()
	handler := &UserHandler{log: log, db: &db.DB{}}

	// Register route
	r.POST("/user", handler.CreateUser)

	// Test case 1: successfully create a user
	t.Run("Success", func(t *testing.T) {
		mockUser := &models.User{Name: "New User", Mail: "new@example.com"}
		mockDB.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()

		// Create request
		userJSON, _ := json.Marshal(mockUser)
		req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Execute request
		r.ServeHTTP(w, req)

		// Validate response
		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID) // Mock assigns ID 1
		assert.Equal(t, mockUser.Name, response.Name)
		assert.Equal(t, mockUser.Mail, response.Mail)

		mockDB.AssertExpectations(t)
	})

	// Test case 2: invalid request body
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

	// Test case 3: missing required fields
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

	// Test case 4: database error
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

// TestUpdateUser exercises the UpdateUser handler
func TestUpdateUser(t *testing.T) {
	r, mockDB, log := setupTest()
	handler := &UserHandler{log: log, db: &db.DB{}}

	// Register route
	r.PUT("/user/:id", handler.UpdateUser)

	// Test case 1: successfully update a user
	t.Run("Success", func(t *testing.T) {
		mockUser := &models.User{ID: 1, Name: "Updated User", Mail: "updated@example.com"}
		mockDB.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()

		// Create request
		userJSON, _ := json.Marshal(mockUser)
		req, _ := http.NewRequest("PUT", "/user/1", bytes.NewBuffer(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Execute request
		r.ServeHTTP(w, req)

		// Validate response
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, mockUser.Name, response.Name)
		assert.Equal(t, mockUser.Mail, response.Mail)

		mockDB.AssertExpectations(t)
	})

	// Test case 2: invalid ID format
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

	// Test case 3: invalid request body
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

	// Test case 4: missing required fields
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

	// Test case 5: database error
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

// TestDeleteUser exercises the DeleteUser handler
func TestDeleteUser(t *testing.T) {
	r, mockDB, log := setupTest()
	handler := &UserHandler{log: log, db: &db.DB{}}

	// Register route
	r.DELETE("/user/:id", handler.DeleteUser)

	// Test case 1: successfully delete a user
	t.Run("Success", func(t *testing.T) {
		mockDB.On("DeleteUser", mock.Anything, uint(1)).Return(nil).Once()

		// Create request
		req, _ := http.NewRequest("DELETE", "/user/1", nil)
		w := httptest.NewRecorder()

		// Execute request
		r.ServeHTTP(w, req)

		// Validate response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User deleted successfully", response["message"])

		mockDB.AssertExpectations(t)
	})

	// Test case 2: invalid ID format
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

	// Test case 3: database error
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
