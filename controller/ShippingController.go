package controllers

import (
	"eCommerce/models"
	"eCommerce/services"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ShippingController struct {
	shippingService *services.ShippingService
}

func NewShippingController(shippingService *services.ShippingService) *ShippingController {
	return &ShippingController{
		shippingService,
	}
}

func (sc *ShippingController) AddShippingAddress(c *gin.Context) {
	var shippingAddress *models.ShippingAddress
	if err := c.ShouldBindJSON(&shippingAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	newShippingAddress, err := sc.shippingService.AddShippingAddress(shippingAddress, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Shipping Address": newShippingAddress})
}

func (sc *ShippingController) DeleteShippingAddress(c *gin.Context) {
	shippingId := c.Query("shippingId")
	if shippingId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing shippingId in query"})
		return
	}

	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	err := sc.shippingService.DeleteShippingAddress(shippingId, userId)
	if err != nil {
		var status int
		var message string

		switch {
		case strings.Contains(err.Error(), "not found"):
			status = http.StatusNotFound
			message = "Shipping address not found"

		case strings.Contains(err.Error(), "unauthorized"):
			status = http.StatusForbidden
			message = "You are not authorized to delete this shipping address"

		case strings.Contains(err.Error(), "ownership"):
			status = http.StatusInternalServerError
			message = "Internal error while checking ownership"

		default:
			status = http.StatusInternalServerError
			message = "Failed to delete shipping address"
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipping address deleted successfully"})
}

func (sc *ShippingController) GetShippingAddresses(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	addresses, err := sc.shippingService.GetShippingAddresses(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"addresses": addresses})
}

func (sc *ShippingController) UpdateShippingAddress(c *gin.Context) {
	var updatedShippingAddress models.ShippingAddress
	if err := c.ShouldBindJSON(&updatedShippingAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shippingId := c.Query("shippingId")
	if shippingId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing shippingId in query"})
		return
	}

	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	shippingAddress, err := sc.shippingService.UpdateShippingAddress(&updatedShippingAddress, userId, shippingId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Updated Address": shippingAddress})
}

func (sc *ShippingController) GetDefaultShippingAddress(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	address, err := sc.shippingService.GetDefaultShippingAddress(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Address": address})
}

func (sc *ShippingController) ChangeDefaultShippingAddress(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	oldDefault := c.Query("oldDefault")
	if oldDefault == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing oldDefault in query"})
		return
	}

	newDefault := c.Query("newDefault")
	if newDefault == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing newDefault in query"})
		return
	}

	fmt.Println(newDefault, oldDefault)

	err := sc.shippingService.ChangeDefaultShippingAddress(userId, oldDefault, newDefault)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated new default shipping address"})
}
