package controllers

import (
	"eCommerce/models"
	"eCommerce/services"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	CartService *services.CartService
}

func NewCartController(cartService *services.CartService) *CartController {
	return &CartController{
		cartService,
	}
}

func (cc *CartController) handleError(err error, c *gin.Context, forbiddenMessage, defaultMessage string) {
	if err != nil {
		var status int
		var message string

		fmt.Println(err.Error())

		switch {
		case strings.Contains(err.Error(), "not found"):
			status = http.StatusNotFound
			message = "Cart not found"

		case strings.Contains(err.Error(), "unauthorized"):
			status = http.StatusForbidden
			message = forbiddenMessage

		case strings.Contains(err.Error(), "ownership"):
			status = http.StatusInternalServerError
			message = "Internal error while checking ownership"

		default:
			status = http.StatusInternalServerError
			message = err.Error()
		}

		c.JSON(status, gin.H{"error": message})
		return
	}
}

func (cc *CartController) DeleteCart(c *gin.Context) {
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

	err := cc.CartService.DeleteCart(cartId, userId)
	if err != nil {
		cc.handleError(err, c, "You are not authorized to delete this cart", "Failed to delete cart")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart deleted successfully"})
}

func (cc *CartController) ViewCart(c *gin.Context) {
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

	cart, err := cc.CartService.ViewCart(cartId, userId)
	if err != nil {
		cc.handleError(err, c, "You are not authorized to view this cart", "Failed to view cart")
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

func (cc *CartController) GetTotalPrice(c *gin.Context) {
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

	totalPrice, err := cc.CartService.GetTotal(cartId, userId)
	if err != nil {
		cc.handleError(err, c, "You are not authorized to check total price of this cart", "Failed to calculate the total price of your cart")
		return
	}

	c.JSON(http.StatusOK, gin.H{"price": totalPrice})
}

func (cc *CartController) DeleteCartItem(c *gin.Context) {
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

	cartItemid := c.Query("cartItemId")
	if cartItemid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing cartItemId in query"})
		return
	}

	err := cc.CartService.DeleteCartItem(cartItemid, userId, cartId)
	if err != nil {

		cc.handleError(err, c, "You are not authorized to delete this item", "Failed to delete item from cart")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item deleted successfully"})
}

func (cc *CartController) EditCartItem(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	var cartItem models.CartItem
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newItem, err := cc.CartService.EditCartItem(&cartItem, userId)
	if err != nil {
		cc.handleError(err, c, "You are not authorized to edit this item", "Failed to edit Item")
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated item": newItem})
}

func (cc *CartController) AddToCart(c *gin.Context) {
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

	var cartItem models.CartItem
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := cc.CartService.AddToCart(cartId, userId, &cartItem)
	if err != nil {
		cc.handleError(err, c, "You are not authorized to add to this cart", "Failed to add to cart")
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

func (cc *CartController) CreateCart(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	err := cc.CartService.CreateCart(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart created"})
}
