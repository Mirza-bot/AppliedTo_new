package routes

import (
	"github.com/gin-gonic/gin"
	"appliedTo/controllers"
)

func SetupJobApplicationRoutes(requireID gin.HandlerFunc) RouteConfig {
	return RouteConfig{
		Prefix: "/job_application",
		Register: func(g *gin.RouterGroup) {
			g.POST("", controllers.CreateJobApplication)
			withID := g.Group("/:id", requireID)
			withID.GET("", controllers.GetJobApplication)
			withID.PUT("", controllers.UpdateJobApplication)
			withID.PATCH("", controllers.PatchJobApplication)
			withID.DELETE("", controllers.DeleteJobApplication)
		},
	}
}
