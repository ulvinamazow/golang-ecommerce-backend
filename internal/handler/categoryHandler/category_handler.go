package categoryhandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	categoryrequestdto "github.com/ulvinamazow/go-backend-project/internal/DTO/request/categoryRequestDTO"
	"github.com/ulvinamazow/go-backend-project/internal/handler"
	"github.com/ulvinamazow/go-backend-project/internal/middleware"
	"github.com/ulvinamazow/go-backend-project/internal/model"
	categoryservice "github.com/ulvinamazow/go-backend-project/internal/service/categoryService"
)

type CategoryHandler struct {
	categoryService categoryservice.CategoryService
}

func NewCategoryController(categoryService categoryservice.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (categoryHandler *CategoryHandler) CreateCategory(c *gin.Context) {
	var req categoryrequestdto.CategoryRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	categoryID, _ := strconv.Atoi(c.Param("id"))

	category := &model.Category{
		ID:   int64(categoryID),
		Name: req.Name,
	}

	if err := categoryHandler.categoryService.Create(c.Request.Context(), category); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (categoryHandler *CategoryHandler) DeleteCategory(c *gin.Context) {
	var req categoryrequestdto.CategoryRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	category, err := categoryHandler.categoryService.GetByID(c.Request.Context(), int64(categoryID))
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	if int64(category.Created_by) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized to delete this category"})
		return
	}

	if err := categoryHandler.categoryService.Delete(c.Request.Context(), int64(categoryID), int64(userID)); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}

func (categoryHandler *CategoryHandler) UpdateCategory(c *gin.Context) {
	var req categoryrequestdto.CategoryRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	category, err := categoryHandler.categoryService.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	if int64(category.Created_by) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized to update this category"})
		return
	}

	category = &model.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := categoryHandler.categoryService.Update(c.Request.Context(), category); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category updated successfully"})
}

func (categoryHandler *CategoryHandler) GetCategoryByID(c *gin.Context) {
	var req categoryrequestdto.CategoryRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	categoryID, _ := strconv.Atoi(c.Param("id"))

	category, err := categoryHandler.categoryService.GetByID(c.Request.Context(), int64(categoryID))
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, category)
}

func (categoryHandler *CategoryHandler) ListCategories(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Param("limit"))
	offset, _ := strconv.Atoi(c.Param("offset"))

	if limit == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	categories, err := categoryHandler.categoryService.List(c.Request.Context(), limit, offset)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, categories)
}
