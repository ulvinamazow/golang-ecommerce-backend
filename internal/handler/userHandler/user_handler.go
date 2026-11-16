package userhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userrequestdto "github.com/ulvinamazow/go-backend-project/internal/DTO/request/userRequestDTO"
	"github.com/ulvinamazow/go-backend-project/internal/handler"
	"github.com/ulvinamazow/go-backend-project/internal/model"
	userservice "github.com/ulvinamazow/go-backend-project/internal/service/userService"
)

type UserHandler struct {
	userService userservice.UserService
}

func NewUserHandler(userService userservice.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (userHandler *UserHandler) CreateUser(c *gin.Context) {
	var req userrequestdto.UserRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": handler.ErrInvalidInput.Error(),
		})
		return
	}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := userHandler.userService.Create(c.Request.Context(), &user); err != nil {
		err := handler.RespondError(c, err)
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
	})
}

func (userHandler *UserHandler) UpdateUser(c *gin.Context) error {
	var req userrequestdto.UserRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": handler.ErrInvalidInput,
		})
		return c.Err()
	}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := userHandler.userService.Update(c.Request.Context(), &user); err != nil {
		// OLD: switch {
		// OLD: }
		return handler.RespondError(c, err)
	}
	return nil
}
