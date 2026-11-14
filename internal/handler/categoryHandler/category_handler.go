package categoryhandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	categoryrequestdto "github.com/ulvinamazow/go-backend-project/internal/DTO/request/categoryRequestDTO"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, _ := strconv.Atoi(c.Param("id"))

	category := &model.Category{
		ID:   int64(productID),
		Name: req.Name,
	}

	if err := categoryHandler.categoryService.Create(c.Request.Context(), category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (categoryHandler *CategoryHandler) DeleteCategory(c *gin.Context) {
	var req categoryrequestdto.CategoryRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	productID, _ := strconv.Atoi(c.Param("id"))

	userID, err := strconv.Atoi(c.GetHeader("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	if err := categoryHandler.categoryService.Delete(c.Request.Context(), int64(productID), int64(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "category can not deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}

func (categoryHandler *CategoryHandler) UpdateCategory(c *gin.Context) {
	var req categoryrequestdto.CategoryRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.GetHeader("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	category := &model.Category{
		ID:          int64(id),
		Name:        req.Name,
		Description: req.Description,
	}

	if err := categoryHandler.categoryService.Update(c.Request.Context(), category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "category can not updated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category updated successfully"})
}

func (categoryHandler *CategoryHandler) GetCategoryByID(c *gin.Context) {
	var req categoryrequestdto.CategoryRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	category, err := categoryHandler.categoryService.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "category can not retrieved"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "categories can not retrieved"})
		return
	}

	c.JSON(http.StatusOK, categories)
}
