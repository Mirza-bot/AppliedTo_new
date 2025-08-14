package routes

import (
	"appliedTo/controllers"
	_ "appliedTo/docs"
	"appliedTo/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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
	exsistingUser := user.Group("/:id", middleware.RequireUserID())
    {
        user.POST("/", controllers.CreateUser)
        exsistingUser.GET("", controllers.GetUser)
        exsistingUser.DELETE("", controllers.DeleteUser)
        exsistingUser.PUT("", controllers.ModifyUser)
        exsistingUser.PATCH("", controllers.ModifyUser)
    }
}

func SetupJobApplicationRoutes(router *gin.Engine) {
	log.Println("Registering job application routes")
	jobApplication := router.Group("/job_application")
	exsistingJobApplication := jobApplication.Group("/:id", middleware.RequireJobApplicationID())
	{
		jobApplication.POST("/", controllers.CreateJobApplication)
		exsistingJobApplication.GET("", controllers.GetJobApplication)
		exsistingJobApplication.PUT("", controllers.UpdateJobApplication)
		exsistingJobApplication.PATCH("", controllers.PatchJobApplication)
		exsistingJobApplication.DELETE("", controllers.DeleteJobApplication)
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

