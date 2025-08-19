package jobapplicationapi

import (
	"appliedTo/internal/platform/http/routes"

	"github.com/gin-gonic/gin"
)

func SetupJobApplicationRoutes(h *Handlers, requireID gin.HandlerFunc) routes.RouteConfig {
    return routes.RouteConfig{
        Prefix: "/job_application",
        Register: func(g *gin.RouterGroup) {
            g.POST("", h.CreateJobApplication)
            withID := g.Group("/:id", requireID)
            withID.GET("", h.GetJobApplication)
            withID.PUT("", h.UpdateJobApplication)
            withID.PATCH("", h.PatchJobApplication)
            withID.DELETE("", h.DeleteJobApplication)
        },
    }
}

