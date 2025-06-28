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

// ShowTask godoc
// @Summary      Получить список задач
// @Description  Возвращает все задачи текущего пользователя
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Success      200 {array} model.Tasks
// @Failure       401 { object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/tasks [get]
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

// GetById godoc
// @Summary      Получить задачу по ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path int true "ID задачи"
// @Success      200 {object} model.Tasks
// @Security     BearerAuth
// @Failure       400 { object} model.ErrorResponse
// @Failure       401 { object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /api/tasks/{id} [get]
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

// DeleteByID godoc
// @Summary      Удалить задачу или пользователя
// @Description  Удаляет задачу или пользователя по ID, в зависимости от параметра choice
// @Tags         tasks
// @Param        id path int true "ID ресурса (task или user)"
// @Param        choice query string true "Тип удаления: task или user"
// @Success      200 {string} string "Успешно удалено"
// @Security     BearerAuth
// @Failure       400 { object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /api/tasks/{id} [delete]
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

// AddTask godoc
// @Summary      Создать задачу
// @Description  Создаёт новую задачу для авторизованного пользователя
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        input body model.Tasks true "Данные задачи"
// @Success      201 {object} model.Tasks
// @Security     BearerAuth
// @Failure       400 { object} model.ErrorResponse
// @Failure       401 { object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /api/tasks [post]
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

// UpdateTask godoc
// @Summary      Обновить задачу
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id    path int          true  "ID задачи"
// @Param        input body model.Tasks   true  "Новые данные"
// @Success      200 {object} model.Tasks
// @Security     BearerAuth
// @Failure       400 { object} model.ErrorResponse
// @Failure       401 { object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /api/tasks/{id} [put]
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
		HandleError(c, errs.ErrNotAccess)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task successfully updated"})
}

// SearchTask  godoc
// @Summary      Поиск задач по подстроке в названии
// @Description  Возвращает список задач, где название содержит указанную подстроку. Поиск зависит от роли пользователя.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        q query string true "Подстрока для поиска в названии задачи"
// @Success      200 {array} model.Tasks
// @Security     BearerAuth
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /api/tasks/ [get]
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

// GetTasksByUserID godoc
// @Summary      Получить задачу через userID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id    path int          true  "ID задачи"
// @Success      200 {object} model.Tasks
// @Security     BearerAuth
// @Failure       400 { object} model.ErrorResponse
// @Failure       401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /api/users/{id}/tasks [get]
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

// ShowFilterTasks godoc
// @Summary      Получить список задач
// @Description  Возвращает все задачи текущего пользователя
// @Tags         tasks
// @Param        q query string true "status tasks:"
// @Accept       json
// @Produce      json
// @Success      200 {array} model.Tasks
// @Failure       401 { object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/filter [get]
func ShowFilterTasks(c *gin.Context) {
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
	query = strings.ToLower(query)
	switch query {
	case "completed":
		{
			tasks, err := service.GetCompletedTasks(role, userID)
			if err != nil {
				HandleError(c, err)
				return
			}
			c.JSON(http.StatusOK, tasks)
		}
	case "incompleted":
		{
			tasks, err := service.GetInCompletedTasks(role, userID)
			if err != nil {
				HandleError(c, err)
				return
			}
			c.JSON(http.StatusOK, tasks)

		}
	case "pending":
		{
			tasks, err := service.GetPendingTasks(role, userID)
			if err != nil {
				HandleError(c, err)
				return
			}
			c.JSON(http.StatusOK, tasks)
		}
	default:
		c.JSON(400, gin.H{"message": "Bad Request"})

	}
}

// FilterByPriority godoc
// @Summary      Получить задачу через userID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Success      200 {object} model.Tasks
// @Security     BearerAuth
// @Failure       400 { object} model.ErrorResponse
// @Failure       401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /api/users/{id}/tasks [get]
func FilterByPriority(c *gin.Context) {
	tasks, err := service.GetTaskByPriority(c.GetString(userRole), c.GetInt(userIDCtx))
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}
