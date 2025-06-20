package controller

import (
	"ToDoList/internal/errs"
	"ToDoList/internal/model"
	"ToDoList/internal/repository"
	"ToDoList/internal/service"
	"ToDoList/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func SignUp(c *gin.Context) {
	var u model.User

	if err := c.ShouldBindJSON(&u); err != nil {
		HandleError(c, err)
		return
	}

	creatorRole := strings.ToLower(c.GetString(userRole))

	if err := service.CreateUser(u, creatorRole); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User successfully registered",
	})
}

func SignIn(c *gin.Context) {
	var u model.UserSignIn
	if err := c.ShouldBindJSON(&u); err != nil {
		HandleError(c, err)
		return
	}
	u.Password = utils.GenerateHash(u.Password)
	user, err := repository.GetUserByUserNameAndPassword(u.Username, u.Password)
	if err != nil {
		HandleError(c, err)
		return
	}
	if user.UserRole == "" {
		user.UserRole = "user"
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.UserRole)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "bearer",
		// можно добавить expires_in, если нужно
	})
}

func GetAllUsers(c *gin.Context) {
	userID := c.GetInt(userIDCtx)
	if userID == 0 {
		HandleError(c, errs.ErrUserNotFound)
		return
	}

	role := strings.ToLower(c.GetString(userRole))
	if role == "" {
		role = "user"
	}

	users, err := service.GetAllUsers(userID, role)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func UpdateUser(c *gin.Context) {
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

	var u model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		HandleError(c, err)
		return
	}

	if err := service.UpdateUser(u, ID, userID, role); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
