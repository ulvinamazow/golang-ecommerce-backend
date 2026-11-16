package userhandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	userrequestdto "github.com/ulvinamazow/go-backend-project/internal/DTO/request/userRequestDTO"
	"github.com/ulvinamazow/go-backend-project/internal/handler"
	"github.com/ulvinamazow/go-backend-project/internal/middleware"
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

func (userHandler *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	middlewareUserID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": handler.ErrUnauthorized.Error()})
		return
	}

	if int64(userID) != middlewareUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": handler.ErrForbidden.Error()})
		return
	}

	var req userrequestdto.UserRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	user := &model.User{
		ID:       int64(userID),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := userHandler.userService.Update(c.Request.Context(), user); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated",
	})
}

func (userHandler *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	middlewareUserID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": handler.ErrUnauthorized.Error()})
		return
	}

	if int64(userID) != middlewareUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": handler.ErrForbidden.Error()})
		return
	}

	if err := userHandler.userService.Delete(c.Request.Context(), int64(userID)); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted",
	})
}

func (userHandler *UserHandler) GetUserByID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	user, err := userHandler.userService.GetByID(c.Request.Context(), int64(userID))
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userID": user.ID,
		"name":   user.Name,
		"email":  user.Email,
	})
}

func (userHandler *UserHandler) GetByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	user, err := userHandler.userService.GetByEmail(c.Request.Context(), email)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userID": user.ID,
		"name":   user.Name,
		"email":  user.Email,
	})
}

func (userHandler *UserHandler) ListUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	users, err := userHandler.userService.List(c.Request.Context(), limit, offset)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}
