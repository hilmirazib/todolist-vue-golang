package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"my-todolist/db"
	"my-todolist/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) func() {
	t.Helper()

	gin.SetMode(gin.TestMode)

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	if err := testDB.AutoMigrate(&models.Todo{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	db.Conn = testDB

	return func() {
		sqlDB, err := testDB.DB()
		if err == nil {
			sqlDB.Close()
		}
		db.Conn = nil
	}
}

func TestListTodos(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	todo := models.Todo{Title: "Sample"}
	if err := db.Conn.Create(&todo).Error; err != nil {
		t.Fatalf("failed to seed todo: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/todos", nil)

	ListTodos(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var todos []models.Todo
	if err := json.Unmarshal(w.Body.Bytes(), &todos); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}

	if todos[0].Title != "Sample" {
		t.Fatalf("expected title 'Sample', got %q", todos[0].Title)
	}
}

func TestCreateTodo(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	body := bytes.NewBufferString(`{"title":"New Todo"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	CreateTodo(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var created models.Todo
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if created.Title != "New Todo" {
		t.Fatalf("expected title 'New Todo', got %q", created.Title)
	}

	if created.ID == 0 {
		t.Fatal("expected todo ID to be set")
	}

	var fromDB models.Todo
	if err := db.Conn.First(&fromDB, created.ID).Error; err != nil {
		t.Fatalf("failed to fetch todo from DB: %v", err)
	}
}

func TestUpdateTodo(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	todo := models.Todo{Title: "Original"}
	if err := db.Conn.Create(&todo).Error; err != nil {
		t.Fatalf("failed to seed todo: %v", err)
	}

	payload := bytes.NewBufferString(`{"title":"Updated","done":true}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodPut, "/todos/"+strconv.Itoa(int(todo.ID)), payload)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(todo.ID))}}

	UpdateTodo(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var updated models.Todo
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if updated.Title != "Updated" {
		t.Fatalf("expected title 'Updated', got %q", updated.Title)
	}

	if !updated.Done {
		t.Fatal("expected todo to be marked as done")
	}
}

func TestUpdateTodoTitleTooLong(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	todo := models.Todo{Title: "Original"}
	if err := db.Conn.Create(&todo).Error; err != nil {
		t.Fatalf("failed to seed todo: %v", err)
	}

	longTitle := strings.Repeat("a", 256)
	bodyBytes, err := json.Marshal(map[string]any{"title": longTitle})
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodPut, "/todos/"+strconv.Itoa(int(todo.ID)), bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(todo.ID))}}

	UpdateTodo(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var fromDB models.Todo
	if err := db.Conn.First(&fromDB, todo.ID).Error; err != nil {
		t.Fatalf("failed to fetch todo from DB: %v", err)
	}

	if fromDB.Title != "Original" {
		t.Fatalf("expected title to remain 'Original', got %q", fromDB.Title)
	}

	if fromDB.Done {
		t.Fatal("expected done status to remain false")
	}
}

func TestDeleteTodo(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	todo := models.Todo{Title: "Delete"}
	if err := db.Conn.Create(&todo).Error; err != nil {
		t.Fatalf("failed to seed todo: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodDelete, "/todos/"+strconv.Itoa(int(todo.ID)), nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(todo.ID))}}

	DeleteTodo(c)
	c.Writer.WriteHeaderNow()

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	var check models.Todo
	err := db.Conn.First(&check, todo.ID).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record to be deleted, got error %v", err)
	}
}
