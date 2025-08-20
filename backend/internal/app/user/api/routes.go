package userapi

import (
	"appliedTo/internal/platform/http/routes"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(h *UserHandlers, requireID gin.HandlerFunc) routes.RouteConfig {
	return routes.RouteConfig{
		Prefix: "/user",
		Register: func(g *gin.RouterGroup) {
			withID := g.Group("/:id", requireID)
			withID.GET("", h.GetUser)
			withID.PUT("", h.UpdateUser)
			withID.PATCH("", h.PatchUser)
			withID.DELETE("", h.DeleteUser)
		},
	}
}

type AdminRouteOpts struct {
	RequireAuth gin.HandlerFunc
	RequireAdmin gin.HandlerFunc
}

func SetupAdminUserRoutes(h *UserHandlers, opts AdminRouteOpts) routes.RouteConfig {
	return routes.RouteConfig{
		Prefix: "/admin/users",
		Register: func(g *gin.RouterGroup) {
			admin := g.Group("", opts.RequireAuth, opts.RequireAdmin)

			admin.POST("", h.CreateUser)
		},
	}

}
