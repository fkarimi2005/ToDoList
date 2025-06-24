package controller

import (
	_ "ToDoList/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RunServer() error {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
