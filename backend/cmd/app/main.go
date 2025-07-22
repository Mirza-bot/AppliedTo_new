package main

import (
	"log"
	"appliedTo/config"
	"appliedTo/routes"

	"github.com/gin-gonic/gin"
)

func main() {
    config.Load()
    config.ConnectDatabase()

    r := gin.Default()

    routes.SetupRoutes(r)

    if err:= r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
