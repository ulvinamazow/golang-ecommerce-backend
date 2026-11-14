package producthandler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	productrequestdto "github.com/ulvinamazow/go-backend-project/internal/DTO/request/productRequestDTO"
	productresponsedto "github.com/ulvinamazow/go-backend-project/internal/DTO/response/productResponseDTO"
	"github.com/ulvinamazow/go-backend-project/internal/handler"
	"github.com/ulvinamazow/go-backend-project/internal/handler/helpers"
	"github.com/ulvinamazow/go-backend-project/internal/middleware"
	"github.com/ulvinamazow/go-backend-project/internal/model"
	productservice "github.com/ulvinamazow/go-backend-project/internal/service/productService"
)

const (
	defaultLimit = 10
	maxLimit     = 100
)

type ProductHandler struct {
	productService productservice.ProductService
}

func NewProductHandler(productService productservice.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// POST /products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": handler.ErrUnauthorized.Error()})
		return
	}

	var req productrequestdto.ProductRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CreatedBy:   userID,
	}

	if err := h.productService.Create(c.Request.Context(), product); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": handler.MsgCreated})
}

// GET /products/:id
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), productID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, toResponse(product))
}

// PUT /products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": handler.ErrUnauthorized.Error()})
		return
	}

	productID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	var req productrequestdto.ProductRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), productID)
	if err != nil {
		respondError(c, err)
		return
	}

	if product.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": handler.ErrForbidden.Error()})
		return
	}

	// Mövcud obyekti dəyişirik ki, CategoryID kimi sahələr itməsin.
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price

	if err := h.productService.Update(c.Request.Context(), product); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": handler.MsgUpdated})
}

// DELETE /products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	productID, err := helpers.GetObjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), productID)
	if err != nil {
		respondError(c, err)
		return
	}

	if product.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": handler.ErrForbidden.Error()})
		return
	}

	if err := h.productService.Delete(c.Request.Context(), productID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": handler.MsgDeleted})
}

// GET /products?limit=10&offset=0
func (h *ProductHandler) ListProducts(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	if err != nil || limit <= 0 || limit > maxLimit {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
		return
	}

	products, err := h.productService.List(c.Request.Context(), limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	// null yox, boş massiv qaytarılsın.
	responses := make([]*productresponsedto.ProductResponseDTO, 0, len(products))
	for _, p := range products {
		responses = append(responses, toResponse(p))
	}

	c.JSON(http.StatusOK, responses)
}

func toResponse(p *model.Product) *productresponsedto.ProductResponseDTO {
	return &productresponsedto.ProductResponseDTO{
		Name:        p.Name,
		CreatedBy:   p.CreatedBy,
		Description: p.Description,
		Price:       p.Price,
		CategoryID:  p.CategoryID,
	}
}
func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, handler.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": handler.ErrNotFound.Error()})
	case errors.Is(err, handler.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": handler.ErrInvalidInput.Error()})
	case errors.Is(err, handler.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": handler.ErrForbidden.Error()})
	default:
		slog.Error("product handler", "path", c.FullPath(), "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": handler.ErrInternal.Error()})
	}
}
