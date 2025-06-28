package controller

import (
	"ToDoList/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userIDCtx           = "userId"
	userRole            = "userRole"
	userNameCtx         = "userName"
)

func checkUserAuthentication(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"error": "Empty authorization header"})
		return
	}
	headerPast := strings.Split(header, " ")
	if len(headerPast) != 2 && headerPast[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"error": "Invalid authorization header"})
		return
	}
	if len(headerPast[1]) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"error": " authorization is empty"})
		return
	}
	accessToken := headerPast[1]
	claims, err := utils.ParseToken(accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"error": err.Error()})
		return
	}

	c.Set(userIDCtx, claims.UserID)
	c.Set(userRole, claims.UserRole)
	c.Set(userNameCtx, claims.Username)
	c.Next()

}
