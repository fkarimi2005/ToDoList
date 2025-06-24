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

// SignUp godoc
// @Summary      Регистрация пользователя
// @Description  Создаёт нового пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body model.User true "Информация о пользователе"
// @Success      201 {object} map[string]string
// @Failure      404 {object} model.ErrorResponse
// @Failure       500 { object} model.ErrorResponse
// @Router       /auth/sign-up [post]
func SignUp(c *gin.Context) {
	var u model.User

	if err := c.ShouldBindJSON(&u); err != nil {
		HandleError(c, err)
		return
	}

	creatorRole := strings.ToLower(c.GetString(userRole))

	if err := service.CreateUser(u, creatorRole); err != nil {
		HandleError(c, errs.ErrUserAlreadyExists)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User successfully registered",
	})
}

// SignIn godoc
// @Summary      Вход в систему
// @Description  Аутентификация пользователя и выдача JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body model.UserSignIn true "Данные пользователя"
// @Success      200 {object} map[string]string
// @Failure      404 {object} model.ErrorResponse
// @Failure       500 { object} model.ErrorResponse
// @Router       /auth/sign-in [post]
func SignIn(c *gin.Context) {
	var u model.UserSignIn
	if err := c.ShouldBindJSON(&u); err != nil {
		HandleError(c, err)
		return
	}
	u.Password = utils.GenerateHash(u.Password)
	user, err := repository.GetUserByUserNameAndPassword(u.Username, u.Password)
	if err != nil {
		HandleError(c, errs.ErrIncorrectLoginOrPassword)
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
	})
}

// GetAllUsers godoc
// @Summary      Получить всех пользователей
// @Description  Возвращает список всех пользователей (только для администратора)
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200 {array} model.User
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      404 {object} model.ErrorResponse
// @Failure       500 { object} model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/users [get]
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

// UpdateUser godoc
// @Summary      Обновить пользователя
// @Description  Обновляет данные пользователя по ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path int true "ID пользователя"
// @Param        input body model.User true "Новые данные пользователя"
// @Success      200 {object} model.User
// @Failure      400 {object} model.ErrorResponse
// @Failure      404 {object} model.ErrorResponse
// @Failure       500 { object} model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/users/{id} [patch]
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
