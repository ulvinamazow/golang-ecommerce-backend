package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/go-backend-project/internal/auth"
	"github.com/ulvinamazow/go-backend-project/internal/config"
	"github.com/ulvinamazow/go-backend-project/internal/database"
	"github.com/ulvinamazow/go-backend-project/internal/middleware"
	"github.com/ulvinamazow/go-backend-project/internal/repository"
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

	connString := config.GetConnectionPool(ctx, *cfg)

	connPool, err := database.Connect(ctx, connString)
	if err != nil {
		panic("database error")
	}

	userRepository := repository.NewPostgresUserRepository(connPool)

	userservice.NewUserService(userRepository)

	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

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
	protected := r.Group("api/auth")

	protected.Use(middleware.AuthMiddleware(autCfg))
	{
		protected.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "hello",
			})
		})
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server can not started: %v", err)
	}

	protected.GET("/login", func(ctx *gin.Context) {
	})

}
