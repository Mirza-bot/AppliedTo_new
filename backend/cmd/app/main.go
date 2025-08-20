package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"appliedTo/docs"
	"appliedTo/internal/app/auth"
	authapi "appliedTo/internal/app/auth/api"
	"appliedTo/internal/app/jobapplication"
	jobapplicationapi "appliedTo/internal/app/jobapplication/api"
	"appliedTo/internal/app/user"
	userapi "appliedTo/internal/app/user/api"
	"appliedTo/internal/platform/config"
	appdb "appliedTo/internal/platform/db"
	"appliedTo/internal/platform/http/middleware"
	"appliedTo/internal/platform/http/routes"
	"appliedTo/internal/platform/security/password"
	"appliedTo/internal/platform/security/token"
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
	jwtIss := &token.JWT{
		Secret:    []byte(cfg.JWTSecret),
		Issuer:    cfg.JWTIssuer,
		AccessTTL: cfg.JWTAccessTTL,
	}

	userService := user.NewService(db, hasher)
	userHandlers := userapi.NewHandlers(userService)

	authService := auth.NewService(db, hasher, jwtIss, userService)
	authHandlers := authapi.NewHandlers(authService)

	jobApplicationService := jobapplication.NewService(db)
	jobApplicationHandlers := jobapplicationapi.NewHandlers(jobApplicationService)

	docs.SwaggerInfo.BasePath = "/api/v1"

	r := gin.Default()
	routes.SetupRoutes(r, "/api/v1",
		authapi.SetupAuthRoutes(authHandlers),
		userapi.SetupUserRoutes(userHandlers, middleware.RequireUserID()),
		jobapplicationapi.SetupJobApplicationRoutes(jobApplicationHandlers, middleware.RequireJobApplicationID()),
	)

	addr := ":" + cfg.AppPort
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}
