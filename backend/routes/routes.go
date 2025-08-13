package routes

import (
	"appliedTo/controllers"
	"log"
    _ "appliedTo/docs"

	"github.com/gin-gonic/gin"
    "github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
)

func SetupAuthRoutes(router *gin.Engine) {
    log.Println("Registering auth routes")
    auth := router.Group("/auth")
    {
        auth.GET("/login", controllers.Login)
        auth.GET("/register", controllers.Register)
    }
}

func SetupUserRoutes(router *gin.Engine) {
	log.Println("Registering user routes")
    user := router.Group("/user")
    {
        user.POST("/", controllers.CreateUser)
        user.GET("/:id", controllers.GetUser)
        user.DELETE("/:id", controllers.DeleteUser)
        user.PATCH("/:id", controllers.ModifyUser)
    }
}

func SetupJobApplicationRoutes(router *gin.Engine) {
	log.Println("Registering job application routes")
	jobApplication := router.Group("/job_application")
	{
		jobApplication.POST("/", controllers.CreateJobApplication)
		jobApplication.GET("/:id", controllers.GetJobApplication)
		jobApplication.PUT("/:id", controllers.UpdateJobApplication)
		jobApplication.PATCH("/:id", controllers.PatchJobApplication)
		jobApplication.DELETE("/:id", controllers.DeleteJobApplication)
	}
}

func SetupRoutes(router *gin.Engine) {
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	log.Println("Setting up routes...")
    SetupAuthRoutes(router)
    SetupUserRoutes(router)
	SetupJobApplicationRoutes(router)
    log.Println("Routes setup completed")
}

