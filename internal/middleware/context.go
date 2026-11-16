package middleware

import "github.com/gin-gonic/gin"

const (
	ContextKeyUserID = "user_id"
	ContextKeyRole   = "role"
)

func GetUserID(c *gin.Context) (int64, error) {
	value, exists := c.Get(ContextKeyUserID)
	if !exists {
		return 0, ErrUserIDNotFound
	}

	userID, ok := value.(int64)
	if !ok {
		return 0, ErrUserIDWrongType
	}

	if userID <= 0 {
		return 0, ErrUserIDEmpty
	}

	return userID, nil
}

func GetRole(c *gin.Context) (string, error) {
	value, exists := c.Get(ContextKeyRole)
	if !exists {
		return "", ErrRoleNotFound
	}

	role, ok := value.(string)
	if !ok {
		return "", ErrRoleWrongType
	}

	if role == "" {
		return "", ErrRoleEmpty
	}

	return role, nil
}
