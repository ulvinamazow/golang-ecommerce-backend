package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/go-backend-project/internal/handler"
)

func GetObjectID(c *gin.Context) (int64, error) {
	objectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || objectID <= 0 {
		return 0, handler.ErrInvalidInput
	}
	return objectID, nil
}
