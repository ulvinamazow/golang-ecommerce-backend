package categoryHandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	categoryrequestdto "github.com/ulvinamazow/go-backend-project/internal/DTO/request/categoryRequestDTO"
	"github.com/ulvinamazow/go-backend-project/internal/model"
	categoryservice "github.com/ulvinamazow/go-backend-project/internal/service/categoryService"
)

type CategoryController struct {
	categoryService categoryservice.CategoryService
}

func NewCategoryController(categoryService categoryservice.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

func (categoryController *CategoryController) CreateCategory(c *gin.Context) {
	var req categoryrequestdto.CategoryRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := &model.Category{
		Name: req.Name,
	}

	if err := categoryController.categoryService.Create(c.Request.Context(), category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "category can not created"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (categoryController *CategoryController) DeleteCategory(c *gin.Context) {
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

	if err := categoryController.categoryService.Delete(c.Request.Context(), int64(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "category can not deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}

func (categoryController *CategoryController) UpdateCategory(c *gin.Context) {
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

	if err := categoryController.categoryService.Update(c.Request.Context(), category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "category can not updated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category updated successfully"})
}

func (categoryController *CategoryController) GetCategoryByID(c *gin.Context) {
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

	category, err := categoryController.categoryService.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "category can not retrieved"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (categoryController *CategoryController) ListCategories(c *gin.Context) {
	limit := 10
	offset := 0

	categories, err := categoryController.categoryService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "categories can not retrieved"})
		return
	}

	c.JSON(http.StatusOK, categories)
}
