package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "appliedTo/docs"
)

type RouteConfig struct {
	Prefix   string
	Use      []gin.HandlerFunc
	Register func(*gin.RouterGroup)
}

func SetupRoutes(r *gin.Engine, basePath string, config ...RouteConfig) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group(basePath)
	log.Println("Setting up routes...")

	for _, conf := range config {
		g := api.Group(conf.Prefix, conf.Use...)
		conf.Register(g)
	}

	log.Println("Routes setup completed")
}
