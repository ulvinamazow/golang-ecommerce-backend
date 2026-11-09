package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/go-backend-project/internal/auth"
	"github.com/ulvinamazow/go-backend-project/internal/config"
	"github.com/ulvinamazow/go-backend-project/internal/middleware"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	config.LoadEnv()

	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	cfg := auth.TokenConfig{
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

	protected.Use(middleware.AuthMiddleware(cfg))
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

}
