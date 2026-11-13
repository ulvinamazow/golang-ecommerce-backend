package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/go-backend-project/internal/auth"
	"github.com/ulvinamazow/go-backend-project/internal/config"
	userrepository "github.com/ulvinamazow/go-backend-project/internal/repository/userRepository"
	userservice "github.com/ulvinamazow/go-backend-project/internal/service/userService"
)

func main() {
	ctx := context.Background()

	gin.SetMode(gin.ReleaseMode)

	if err := config.LoadEnv(); err != nil {
		panic("environment variables can not load")
	}

	cfg, err := config.Load()
	if err != nil {
		panic("error")
	}

	connPool, err := config.GetConnectionPool(ctx, *cfg)

	if err != nil {
		panic("error")
	}

	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	userRepository := userrepository.NewPostgresUserRepository(connPool)

	userservice.NewUserService(userRepository)

	autCfg := auth.TokenConfig{
		SecretKey: []byte(os.Getenv("JWT_SECRET_KEY")),
		Issuer:    "ecommerce-api",
		Audience:  "ecommerce-client",
		TTL:       20 * time.Minute,
	}

	/*
		public := r.Group("/public")
		{
			public.POST("/login", authHandler.Login)
			public.POST("/register", authHandler.Register)
			public.POST("/refresh", authHandler.Refresh)
		}

	*/

}
