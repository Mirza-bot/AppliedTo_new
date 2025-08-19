package userapi

import (
	"appliedTo/internal/platform/http/routes"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(h *UserHandlers, requireID gin.HandlerFunc) routes.RouteConfig {
	return routes.RouteConfig{
		Prefix: "/user",
		Register: func(g *gin.RouterGroup) {
			g.POST("", h.CreateUser)

			withID := g.Group("/:id", requireID)
			withID.GET("", h.GetUser)
			withID.PUT("", h.UpdateUser)
			withID.PATCH("", h.PatchUser)
			withID.DELETE("", h.DeleteUser)
		},
	}
}
