package productHandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	productrequestdto "github.com/ulvinamazow/go-backend-project/internal/DTO/request/productRequestDTO"
	productresponsedto "github.com/ulvinamazow/go-backend-project/internal/DTO/response/productResponseDTO"
	"github.com/ulvinamazow/go-backend-project/internal/handler"
	"github.com/ulvinamazow/go-backend-project/internal/model"
	productservice "github.com/ulvinamazow/go-backend-project/internal/service/productService"
)

type ProductController struct {
	productService productservice.ProductService
}

func NewProductController(productService productservice.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

func (productController *ProductController) CreateProduct(c *gin.Context) {
	var req productrequestdto.ProductRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": handler.ErrInvalidInput,
		})
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": handler.ErrForbidden})
		return
	}

	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CreatedBy:   userID,
	}

	if err := productController.productService.Create(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": handler.ErrInvalidInput,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": handler.StatusCreated,
	})
}

func (productController *ProductController) GetProductByID(c *gin.Context) (*productresponsedto.ProductResponseDTO, error) {
	// var req productrequestdto.ProductRequestDTO
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "invalid input",
	// 	})
	// 	return nil, err
	// }

	productID, err := getObjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": handler.ErrInvalidInput,
		})
		return nil, handler.ErrInvalidInput
	}

	product, err := productController.productService.GetByID(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": handler.ErrInvalidInput,
		})
		return nil, handler.ErrInvalidInput
	}

	return &productresponsedto.ProductResponseDTO{
		Name:        product.Name,
		CreatedBy:   product.CreatedBy,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
	}, nil
}

func (productController *ProductController) UpdateProduct(c *gin.Context) {
	productID, err := getObjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": handler.ErrInvalidInput,
		})
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": handler.ErrForbidden})
		return
	}

	var req productrequestdto.ProductRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": handler.ErrInvalidInput,
		})
		return
	}

	product, err := productController.productService.GetByID(c.Request.Context(), productID)
	if userID != product.CreatedBy {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": handler.ErrUnauthorized,
		})
		return
	}

	product = &model.Product{
		ID:          productID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CreatedBy:   userID,
	}

	if err := productController.productService.Update(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": handler.ErrInvalidInput,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": handler.StatusCreated,
	})

}

func (productController *ProductController) DeleteProduct(c *gin.Context) {
	productID, err := getObjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": handler.ErrInvalidInput,
		})
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": handler.ErrForbidden})
		return
	}

	product, err := productController.productService.GetByID(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": handler.ErrInvalidInput,
		})
		return
	}

	if product.CreatedBy != userID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": handler.ErrUnauthorized,
		})
		return
	}

	err = productController.productService.Delete(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": handler.ErrInvalidInput,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": handler.StatusDeleted,
	})
}

func (productController *ProductController) ListProducts(c *gin.Context) ([]*productresponsedto.ProductResponseDTO, error) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	products, err := productController.productService.List(c.Request.Context(), limit, offset)
	if err != nil {
		return nil, handler.ErrInvalidInput
	}

	var productResponses []*productresponsedto.ProductResponseDTO
	for _, product := range products {
		productResponse := &productresponsedto.ProductResponseDTO{
			Name:        product.Name,
			CreatedBy:   product.CreatedBy,
			Description: product.Description,
			Price:       product.Price,
			CategoryID:  product.CategoryID,
		}
		productResponses = append(productResponses, productResponse)
	}
	return productResponses, nil
}

/*
func (productController *ProductController) GetProductsByCategoryID(c *gin.Context) ([]*productresponsedto.ProductResponseDTO, error) {
	idStr := c.Param("id")
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return nil, err
	}

	products, err := productController.productService.GetByID(c.Request.Context(), int64(categoryID))
	if err != nil {
		return nil, err
	}

	var productResponses []*productresponsedto.ProductResponseDTO
	for _, product := range products {
		productResponse := &productresponsedto.ProductResponseDTO{
			Name:        products.Name,
			CreatedBy:   products.CreatedBy,
			Description: products.Description,
			Price:       products.Price,
			CategoryID:  products.CategoryID,
		}
		productResponses = append(productResponses, productResponse)
	}
	return productResponses, nil
}

func (productController *ProductController) GetProductsByUserID(c *gin.Context) ([]*productresponsedto.ProductResponseDTO, error) {
	// userID, err := getUserID(c)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorization error"})
	// 	return []*productresponsedto.ProductResponseDTO{}, err
	// }

	productID, err := getObjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return []*productresponsedto.ProductResponseDTO{}, handler.ErrInvalidInput
	}

	productController.productService.GetByID(c.Request.Context(), productID)

	var productResponses []*productresponsedto.ProductResponseDTO
	for _, product := range products {
		productResponse := &productresponsedto.ProductResponseDTO{
			Name:        product.Name,
			CreatedBy:   product.CreatedBy,
			Description: product.Description,
			Price:       product.Price,
			CategoryID:  product.CategoryID,
		}
		productResponses = append(productResponses, productResponse)
	}
	return productResponses, nil
}
*/

func getUserID(c *gin.Context) (int64, error) {
	userIdStr := c.GetHeader("user_id")

	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		return 0, err
	}

	return int64(userID), nil
}

func getObjectID(c *gin.Context) (int64, error) {
	idStr := c.Param("id")

	objectID, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return int64(objectID), nil
}
