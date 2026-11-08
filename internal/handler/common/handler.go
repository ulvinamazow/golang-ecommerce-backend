package common

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleUserRequest(c *gin.Context) {
	router := gin.Default()

	router.GET("homepage", handleHomepage)
}

func handleHomepage(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		log.Fatal("Use the GET method!")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, welcome to my site",
	})

}
