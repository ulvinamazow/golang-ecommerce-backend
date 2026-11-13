package productHandler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	productrequestdto "github.com/ulvinamazow/go-backend-project/internal/DTO/request/productRequestDTO"
	productresponsedto "github.com/ulvinamazow/go-backend-project/internal/DTO/response/productResponseDTO"
	"github.com/ulvinamazow/go-backend-project/internal/model"
	productservice "github.com/ulvinamazow/go-backend-project/internal/service/productService"
)

type ProductController struct {
	productService productservice.ProductService
}

func NewProductController(productService productservice.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

func (productcontroller *ProductController) CreateProduct(
	c *gin.Context,
	productRequest *productrequestdto.ProductRequestDTO) {

	id, err := strconv.Atoi(c.GetHeader("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	product := &model.Product{
		ID:          int64(id),
		Name:        productRequest.Name,
		Description: productRequest.Description,
		Price:       productRequest.Price,
	}

	if err := productcontroller.productService.Create(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product created",
	})
	return
}

func (pc *ProductController) GetProductByID(ctx context.Context, id int64) (*productresponsedto.ProductResponseDTO, error) {
	product, err := pc.productService.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &productresponsedto.ProductResponseDTO{
		Name:        product.Name,
		CreatedBy:   product.CreatedBy,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
	}, nil
}

func (productcontroller *ProductController) UpdateProduct(c *gin.Context, productRequest *productrequestdto.ProductRequestDTO) *productresponsedto.ProductResponseDTO {
	id, err := strconv.Atoi(c.GetHeader("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return nil
	}

	product := &model.Product{
		ID:          productRequest.ID,
		Name:        productRequest.Name,
		Description: productRequest.Description,
		Price:       productRequest.Price,
		CreatedBy:   int64(id),
	}

	updatedProduct, err := productcontroller.productService.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	return &productresponsedto.ProductResponseDTO{
		Name:        updatedProduct.Name,
		CreatedBy:   updatedProduct.CreatedBy,
		Description: updatedProduct.Description,
		Price:       updatedProduct.Price,
		CategoryID:  updatedProduct.CategoryID,
	}, nil
}

func (pc *ProductController) DeleteProduct(ctx context.Context, id int64) error {
	err := pc.productService.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (pc *ProductController) ListProducts(ctx context.Context, limit, offset int) ([]*productresponsedto.ProductResponseDTO, error) {
	products, err := pc.productService.List(ctx, limit, offset)
	if err != nil {
		return nil, err
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

func (pc *ProductController) GetProductsByCategoryID(ctx context.Context, categoryID int64) ([]*productresponsedto.ProductResponseDTO, error) {
	products, err := pc.productService.GetByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
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

func (pc *ProductController) GetProductsByUserID(ctx context.Context, userID int64) ([]*productresponsedto.ProductResponseDTO, error) {
	products, err := pc.productService.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
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
