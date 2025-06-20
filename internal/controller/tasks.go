package controller

import (
	"ToDoList/internal/errs"
	"ToDoList/internal/model"
	"ToDoList/internal/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "server up and running",
	})
}

func ShowTask(c *gin.Context) {
	userID := c.GetInt(userIDCtx)
	if userID == 0 {
		HandleError(c, errs.ErrUserNotFound)
		return
	}
	role := strings.ToLower(c.GetString(userRole))
	if role == "" {
		role = "user"
	}

	tasks, err := service.ShowTask(role, userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func GetById(c *gin.Context) {
	IDStr := c.Param("id")
	ID, err := strconv.Atoi(IDStr)
	if err != nil || ID < 0 {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	userID := c.GetInt(userIDCtx)
	if userID == 0 {
		HandleError(c, errs.ErrUserNotFound)
		return
	}

	role := strings.ToLower(c.GetString(userRole))
	if role == "" {
		role = "user"
	}

	task, err := service.GetTaskByID(ID, userID, role)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func DeleteByID(c *gin.Context) {
	userID := c.GetInt(userIDCtx)
	if userID == 0 {
		HandleError(c, errs.ErrUserNotFound)
		return
	}

	role := strings.ToLower(c.GetString(userRole))
	if role == "" {
		role = "user"
	}

	IDStr := c.Param("id")
	ID, err := strconv.Atoi(IDStr)
	if err != nil || ID < 0 {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	queryChoice := strings.ToLower(c.Query("choice"))
	if queryChoice == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query parameter 'choice' is required (user or task)",
		})
		return
	}

	switch queryChoice {
	case "user":
		err := service.DeleteUsers(ID, userID, role)
		if err != nil {
			fmt.Printf("DeleteUsers error: %v\n", err) // debug log
			HandleError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User has been deleted"})
	case "task":
		err := service.DeleteTask(ID, userID, role)
		if err != nil {
			HandleError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Task has been successfully deleted"})
	default:
		HandleError(c, errs.ErrInvalidOperationType)
	}
}

func AddTask(c *gin.Context) {
	userID := c.GetInt(userIDCtx)
	if userID == 0 {
		HandleError(c, errs.ErrUserNotFound)
		return
	}

	role := strings.ToLower(c.GetString(userRole))
	if role == "" {
		role = "user"
	}

	var task model.Tasks
	if err := c.ShouldBindJSON(&task); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	if err := service.CreateTask(task, role, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Task successfully added"})
}

func UpdateTask(c *gin.Context) {
	IDStr := c.Param("id")
	ID, err := strconv.Atoi(IDStr)
	if err != nil || ID < 0 {
		HandleError(c, errors.New("Invalid ID format. Must be a non-negative integer."))
		return
	}

	userID := c.GetInt(userIDCtx)
	if userID == 0 {
		HandleError(c, errs.ErrUserNotFound)
		return
	}

	role := strings.ToLower(c.GetString(userRole))
	if role == "" {
		role = "user"
	}

	var d model.DoneTasks
	if err := c.ShouldBindJSON(&d); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	if err := service.UpdateTask(d, ID, userID, role); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task successfully updated"})
}

func SearchTask(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Query parameter 'q' is missing"})
		return
	}

	userID := c.GetInt(userIDCtx)
	if userID == 0 {
		HandleError(c, errs.ErrUserNotFound)
		return
	}

	role := strings.ToLower(c.GetString(userRole))
	if role == "" {
		role = "user"
	}

	tasks, err := service.SearchTask(query, role, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func GetTasksByUserID(c *gin.Context) {
	userID := c.GetInt(userIDCtx)
	if userID == 0 {
		HandleError(c, errs.ErrUserNotFound)
		return
	}

	role := strings.ToLower(c.GetString(userRole))
	if role == "" {
		role = "user"
	}

	IDStr := c.Param("id")
	requestID, err := strconv.Atoi(IDStr)
	if err != nil || requestID < 0 {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	tasks, err := service.GetTasksByUserID(requestID, userID, role)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}
