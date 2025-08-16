package controllers

import (
	"eCommerce/models"
	"eCommerce/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type VendorController struct {
	vendorService *services.VendorService
}

func NewVendorController(vendorService *services.VendorService) *VendorController {
	return &VendorController{
		vendorService,
	}
}

func (vc *VendorController) AddProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vendorIdRaw, exists := c.Get("UserId")
	vendorId, ok := vendorIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	newProduct, err := vc.vendorService.AddProduct(&product, vendorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": newProduct})
}

func (vc *VendorController) GetProductById(c *gin.Context) {
	productId := c.Query("productId")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
		return
	}
	vendorIdRaw, exists := c.Get("UserId")
	vendorId, ok := vendorIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	product, err := vc.vendorService.GetProductByID(productId, vendorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (vc *VendorController) GetVendorProducts(c *gin.Context) {
	vendorIdRaw, exists := c.Get("UserId")
	vendorId, ok := vendorIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	var products []*models.Product
	products, err := vc.vendorService.GetVendorProducts(vendorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (vc *VendorController) DeleteProduct(c *gin.Context) {
	productId := c.Query("productId")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
		return
	}
	vendorIdRaw, exists := c.Get("UserId")
	vendorId, ok := vendorIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	err := vc.vendorService.DeleteProduct(productId, vendorId)
	if err != nil {
		var status int
		var message string

		switch {
		case strings.Contains(err.Error(), "not found"):
			status = http.StatusNotFound
			message = "product not found"

		case strings.Contains(err.Error(), "unauthorized"):
			status = http.StatusForbidden
			message = "You are not authorized to delete this product"

		case strings.Contains(err.Error(), "ownership"):
			status = http.StatusInternalServerError
			message = "Internal error while checking ownership"

		default:
			status = http.StatusInternalServerError
			message = "Failed to delete product"
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (vc *VendorController) UpdateProduct(c *gin.Context) {
	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updatedProduct.ProductId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ProductId is required"})
		return
	}

	vendorIdRaw, exists := c.Get("UserId")
	vendorId, ok := vendorIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	product, err := vc.vendorService.UpdateProduct(&updatedProduct, vendorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Updated Product": product})
}
