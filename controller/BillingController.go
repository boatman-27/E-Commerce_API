package controllers

import (
	"eCommerce/models"
	"eCommerce/services"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BillingController struct {
	billingService *services.BillingService
}

func NewBillingController(billingService *services.BillingService) *BillingController {
	return &BillingController{
		billingService,
	}
}

func (bc *BillingController) AddBillingAddress(c *gin.Context) {
	var billingAddress models.BillingAddress
	if err := c.ShouldBindJSON(&billingAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	newBillingAddress, err := bc.billingService.AddBillingAddress(&billingAddress, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Billing Address": newBillingAddress})
}

func (bc *BillingController) DeleteBillingAddress(c *gin.Context) {
	billingId := c.Query("billingId")
	if billingId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing billingId in query"})
		return
	}

	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	err := bc.billingService.DeleteBillingAddress(billingId, userId)
	if err != nil {
		var status int
		var message string

		switch {
		case strings.Contains(err.Error(), "not found"):
			status = http.StatusNotFound
			message = "Billing address not found"

		case strings.Contains(err.Error(), "unauthorized"):
			status = http.StatusForbidden
			message = "You are not authorized to delete this billing address"

		case strings.Contains(err.Error(), "ownership"):
			status = http.StatusInternalServerError
			message = "Internal error while checking ownership"

		default:
			status = http.StatusInternalServerError
			message = "Failed to delete billing address"
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Billing address deleted successfully"})
}

func (bc *BillingController) GetBillingAddresses(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	addresses, err := bc.billingService.GetBillingAddresses(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"addresses": addresses})
}

func (bc *BillingController) UpdateBillingAddress(c *gin.Context) {
	var updatedBillingAddress models.BillingAddress
	if err := c.ShouldBindJSON(&updatedBillingAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	billingId := c.Query("billingId")
	if billingId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing billingId in query"})
		return
	}

	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	billingAddress, err := bc.billingService.UpdateBillingAddress(&updatedBillingAddress, userId, billingId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Updated Address": billingAddress})
}

func (bc *BillingController) GetDefaultBillingAddress(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	address, err := bc.billingService.GetDefaultBillingAddress(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Address": address})
}

func (bc *BillingController) ChangeDefaultBillingAddress(c *gin.Context) {
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

	err := bc.billingService.ChangeDefaultBillingAddress(userId, oldDefault, newDefault)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated new default billing address"})
}
