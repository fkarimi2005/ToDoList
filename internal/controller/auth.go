package controller

import (
	"ToDoList/errs"
	"ToDoList/model"
	"ToDoList/repository"
	"ToDoList/service"
	"ToDoList/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func SignUp(c *gin.Context) {
	var u model.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		HandleError(c, err)
		return
	}

	role := c.GetString(userRole)
	if role == "" {
		c.JSON(200, gin.H{
			"message": " role is empty",
		})
		return
	}
	role = strings.ToLower(role)
	err = service.CreateUser(u, role)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, gin.H{
		"message": "success addet user",
	})

}
func SignIn(c *gin.Context) {
	var u model.UserSignIn
	err := c.ShouldBindJSON(&u)
	if err != nil {
		HandleError(c, err)
		return
	}
	user, err := repository.GetUserByUserNameAndPassword(u.Username, u.Password)
	if err != nil {
		HandleError(c, err)
		return
	}
	token, err := utils.GenerateToken(user.ID, user.Username, user.UserRole)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "bearer",
		// если хочешь указать время жизни в секундах
	})
}

func GetAllUsers(c *gin.Context) {
	userID := c.GetInt(userIDCtx)
	if userID == 0 {
		HandleError(c, errs.ErrUserNotFound)
		return
	}
	role := c.GetString(userRole)
	if role == "" {
		c.JSON(200, gin.H{
			"message": " role is empty",
		})
		return
	}
	role = strings.ToLower(role)
	users, err := service.GetAllUsers(userID, role)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, users)

}
func UpdateUser(c *gin.Context) {
	IDStr := c.Param("id")
	ID, err := strconv.Atoi(IDStr)
	if ID < 0 || err != nil {
		HandleError(c, err)
		return
	}
	userID := c.GetInt(userIDCtx)
	if userID == 0 {
		HandleError(c, errs.ErrUserNotFound)
		return
	}
	role := c.GetString(userRole)
	if role == "" {
		c.JSON(200, gin.H{
			"message": " role is empty",
		})
		return
	}
	role = strings.ToLower(role)
	var u model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		HandleError(c, err)
		return
	}
	if err := service.UpdateUser(u, ID, userID, role); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, gin.H{
		"message": "User updated successfully",
	})

}

//func ChangeRole(c *gin.Context) {
//	role := c.Query("role")
//	if role == "" {
//		c.JSON(200, gin.H{
//			"message": "Query role is empty",
//		})
//	}
//	var u model.User
//	if err := c.ShouldBindJSON(&u); err != nil {
//		HandleError(c, err)
//		return
//	}
//	err := service.ChangeRoleUser(u.Username, role)
//	if err != nil {
//		HandleError(c, err)
//		return
//	}
//	c.JSON(201, gin.H{
//		"message": "success change role user",
//	})
//}
