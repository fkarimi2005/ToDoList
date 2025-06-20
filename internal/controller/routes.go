package controller

import (
	"github.com/gin-gonic/gin"
)

func RunServer() error {
	router := gin.Default()
	router.GET("/", Ping)
	authG := router.Group("/auth")
	{
		authG.POST("/sign-up", SignUp)
		authG.POST("/sign-in", SignIn)

	}
	apiG := router.Group("/api", checkUserAuthentication)
	userG := apiG.Group("/users")
	{
		userG.GET("", GetAllUsers)
		userG.PATCH("/:id", UpdateUser)
		userG.GET("/:id/tasks", GetTasksByUserID)

	}

	taskG := apiG.Group("/tasks")

	{

		taskG.GET("", ShowTask)
		taskG.GET("/", SearchTask)
		taskG.GET("/:id", GetById)
		taskG.DELETE("/:id/", DeleteByID)
		//taskG.PATCH("/:id", UpStatusTask)
		taskG.POST("", AddTask)
		taskG.PUT("/:id", UpdateTask)

	}
	return router.Run(":8089")

}
