package category

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/go-backend-project/internal/model"
	"github.com/ulvinamazow/go-backend-project/internal/service"
)

type CategoryController struct {
	categoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

func (categoryController *CategoryController) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest

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
	
}