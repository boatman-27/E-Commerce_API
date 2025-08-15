package controllers

import (
	"eCommerce/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CheckoutController struct {
	checkoutService *services.CheckoutService
}

func NewCheckoutController(checkoutService *services.CheckoutService) *CheckoutController {
	return &CheckoutController{
		checkoutService: checkoutService,
	}
}

func (chc *CheckoutController) OrderSummary(c *gin.Context) {
	cartId := c.Query("cartId")
	if cartId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing cartId in query"})
		return
	}

	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	summary, err := chc.checkoutService.OrderSummary(cartId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"summary": summary})
}

func (chc *CheckoutController) ConfrimPurchase(c *gin.Context) {
	cartId := c.Query("cartId")
	if cartId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing cartId in query"})
		return
	}

	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	emailRaw, exists := c.Get("Email")
	email, ok := emailRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	summary, err := chc.checkoutService.ConfirmPurchase(cartId, userId, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"summary": summary})
}
