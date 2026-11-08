package middleware

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func FileLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		f, _ := os.OpenFile("gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

		r := gin.New()
		r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			Output: io.MultiWriter(f, os.Stdout),
		}))
	}
}
