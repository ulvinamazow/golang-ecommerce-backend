package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/go-backend-project/internal/auth"
	"github.com/ulvinamazow/go-backend-project/internal/config"
	producthandler "github.com/ulvinamazow/go-backend-project/internal/handler/productHandler"
	"github.com/ulvinamazow/go-backend-project/internal/middleware"
	productrepository "github.com/ulvinamazow/go-backend-project/internal/repository/productRepository"
	productservice "github.com/ulvinamazow/go-backend-project/internal/service/productService"
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

	// userRepository := userrepository.NewPostgresUserRepository(connPool)
	productRepository := productrepository.NewPostgresProductRepository(connPool)
	// categoryRepository := categoryrepository.NewPostgresCategoryRepository(connPool)

	// userService := userservice.NewUserService(userRepository)
	productService := productservice.NewProductService(productRepository)
	// categoryService := categoryservice.NewCategoryService(categoryRepository)

	// userHandler := userhandler.NewUserHandler(userService)
	productHandler := producthandler.NewProductHandler(productService)
	// categoryHandler := categoryhandler.NewCategoryController(categoryService)

	autCfg := auth.TokenConfig{
		SecretKey: []byte(os.Getenv("JWT_SECRET_KEY")),
		Issuer:    "ecommerce-api",
		Audience:  "ecommerce-client",
		TTL:       20 * time.Minute,
	}

	r.Use(middleware.AuthMiddleware(autCfg))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	r.GET("/", productHandler.ListProducts)
	r.POST("/products", productHandler.CreateProduct)
	r.PUT("/products/:id", productHandler.UpdateProduct)
	r.DELETE("/products/:id", productHandler.DeleteProduct)
	r.GET("/products/:id", productHandler.GetProductByID)

	r.Run(":8080")

}
