package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"appliedTo/docs"
	"appliedTo/handlers"
	"appliedTo/internal/config"
	appdb "appliedTo/internal/db"
	"appliedTo/internal/routes"
	"appliedTo/internal/security/password"
	userservice "appliedTo/internal/services/user_service"
	"appliedTo/middleware"
)

func main() {
	cfg := config.Load()

	db, err := appdb.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	if err := appdb.Migrate(db); err != nil {
		log.Fatalf("db migrate: %v", err)
	}

	hasher := password.NewBcrypt(password.WithCost(cfg.BcryptCost))

	userSvc := userservice.NewUserService(db, hasher)
	userHandlers := handlers.NewUserHandlers(userSvc)

	docs.SwaggerInfo.BasePath = "/api/v1"

	r := gin.Default()
	routes.SetupRoutes(r, "/api/v1",
	routes.SetupAuthRoutes(),
	routes.SetupUserRoutes(userHandlers, middleware.RequireUserID()),
	routes.SetupJobApplicationRoutes(middleware.RequireJobApplicationID()),
	)


	addr := ":" + cfg.AppPort
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}

