package authapi

import (
    "appliedTo/internal/platform/http/routes"
    "github.com/gin-gonic/gin"
)

func SetupAuthRoutes(h *Handlers) routes.RouteConfig {
    return routes.RouteConfig{
        Prefix: "/auth",
        Register: func(g *gin.RouterGroup) {
            g.POST("/login", h.Login)
            g.POST("/register", h.Register)
        },
    }
}

