package routes

import (
	"github.com/gin-gonic/gin"
	"appliedTo/handlers"
)

func SetupUserRoutes(h *handlers.UserHandlers, requireID gin.HandlerFunc) RouteConfig {
	return RouteConfig{
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
