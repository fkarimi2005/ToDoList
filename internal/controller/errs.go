package controller

import (
	"ToDoList/internal/errs"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	switch {
	case err == nil:
		return

	case errors.Is(err, errs.ErrValidationFailed),
		errors.Is(err, errs.ErrAlreadyDeleted),
		errors.Is(err, errs.ErrUserAlreadyExists),
		errors.Is(err, errs.ErrInvalidID):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	case errors.Is(err, errs.ErrIncorrectLoginOrPassword):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return

	case errors.Is(err, errs.ErrNoPermissionsToCreateTask):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return

	case errors.Is(err, errs.ErrNotFoundID),
		errors.Is(err, errs.ErrUserNotFound),
		errors.Is(err, errs.ErrTaskNotFound),
		errors.Is(err, errs.ErrNotFoud):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return

	case errors.Is(err, errs.ErrSomethingWentWrong):
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error: " + err.Error()})
		return
	}

}
