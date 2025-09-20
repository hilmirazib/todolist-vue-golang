package handlers

import (
	"my-todolist/db"
	"my-todolist/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type CreateTodoDTO struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}
type UpdateTodoDTO struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

func ListTodos(c *gin.Context) {
	var todos []models.Todo
	if err := db.Conn.Order("id DESC").Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB_ERROR", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func CreateTodo(c *gin.Context) {
	var body CreateTodoDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR"})
		return
	}
	if err := validate.Struct(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "detail": err.Error()})
		return
	}
	t := models.Todo{Title: body.Title}
	if err := db.Conn.Create(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB_ERROR", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func UpdateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_ID"})
		return
	}
	var body UpdateTodoDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR"})
		return
	}
	var t models.Todo
	if err := db.Conn.First(&t, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NOT_FOUND"})
		return
	}
	if body.Title != nil {
		t.Title = *body.Title
	}
	if body.Done != nil {
		t.Done = *body.Done
	}
	if err := db.Conn.Save(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB_ERROR"})
		return
	}
	c.JSON(http.StatusOK, t)
}

func DeleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_ID"})
		return
	}
	if err := db.Conn.Delete(&models.Todo{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB_ERROR"})
		return
	}
	c.Status(http.StatusNoContent)
}
