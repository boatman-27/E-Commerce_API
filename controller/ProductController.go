package controllers

import (
	"eCommerce/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService *services.ProductService
}

func NewProductController(productService *services.ProductService) *ProductController {
	return &ProductController{
		productService,
	}
}

func (pc *ProductController) GetProducts(c *gin.Context) {
	filters := map[string]string{}
	if value := c.Query("category"); value != "" {
		filters["category"] = value
	}
	if value := c.Query("brand"); value != "" {
		filters["brand"] = value
	}
	if value := c.Query("min_price"); value != "" {
		filters["min_price"] = value
	}
	if value := c.Query("max_price"); value != "" {
		filters["max_price"] = value
	}
	if value := c.Query("discount"); value != "" {
		filters["discount"] = value
	}
	if value := c.Query("keyword"); value != "" {
		filters["keyword"] = value
	}
	if value := c.Query("sort_by"); value != "" {
		filters["sort_by"] = value
	}
	if value := c.Query("order"); value != "" {
		filters["order"] = value
	}
	if value := c.Query("limit"); value != "" {
		filters["limit"] = value
	}
	if value := c.Query("offset"); value != "" {
		filters["offset"] = value
	}

	products, err := pc.productService.GetProducts(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}
