package middleware

import "github.com/gin-gonic/gin"

const (
	ContextKeyUserID = "user_id"
	ContextKeyRole   = "role"
)

func GetUserID(c *gin.Context) (string, error) {
	value, exists := c.Get(ContextKeyUserID)

	if !exists {
		return "", ErrUserIDNotFound
	}

	userID, ok := value.(string)

	if !ok {
		return "", ErrUserIDWrongType
	}

	if value == "" {
		return "", ErrUserIDEmpty
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
		return "", ErrRoleNotFound
	}

	return role, nil
}
