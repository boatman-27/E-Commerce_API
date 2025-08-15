package controllers

import (
	"eCommerce/models"
	"eCommerce/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WishlistController struct {
	WishlistService *services.WishlistService
}

func NewWishlistController(wishlistService *services.WishlistService) *WishlistController {
	return &WishlistController{
		wishlistService,
	}
}

func (wc *WishlistController) CreateWishlist(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	var info *models.WishlistInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := wc.WishlistService.CreateWishlist(userId, info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "wishlist created successfully"})
}

func (wc *WishlistController) DeleteWishlist(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	wishlistId := c.Query("wishlistId")
	if wishlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wishlistId is missing in Query"})
		return
	}

	err := wc.WishlistService.DeleteWishlist(userId, wishlistId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "wishlist deleted successfully"})
}

func (wc *WishlistController) EditWishlist(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	wishlistId := c.Query("wishlistId")
	if wishlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wishlistId is missing in Query"})
		return
	}

	var info *models.WishlistInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := wc.WishlistService.EditWishlist(userId, wishlistId, info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "wishlist edited successfully"})
}

func (wc *WishlistController) DeleteWishlistItem(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	wishlistId := c.Query("wishlistId")
	if wishlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wishlistId is missing in Query"})
		return
	}

	wishlistItemId := c.Query("wishlistItemId")
	if wishlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wishlistItemId is missing in Query"})
		return
	}

	err := wc.WishlistService.DeleteWishlistItem(userId, wishlistId, wishlistItemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "wishlist item deleted successfully"})
}

func (wc *WishlistController) MoveToCart(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	wishlistId := c.Query("wishlistId")
	if wishlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wishlistId is missing in Query"})
		return
	}

	err := wc.WishlistService.MoveToCart(userId, wishlistId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "wishlist items moved to cart"})
}

func (wc *WishlistController) AddToWishlist(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	wishlistId := c.Query("wishlistId")
	if wishlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wishlistId is missing in Query"})
		return
	}

	productId := c.Query("productId")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId is missing in Query"})
		return
	}

	wishlist, err := wc.WishlistService.AddToWishlist(userId, wishlistId, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, gin.H{"wishlist": wishlist})
}

func (wc *WishlistController) MoveItemToCart(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	wishlistId := c.Query("wishlistId")
	if wishlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wishlistId is missing in Query"})
		return
	}

	wishlistItemId := c.Query("wishlistItemId")
	if wishlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wishlistItemId is missing in Query"})
		return
	}

	err := wc.WishlistService.MoveItemToCart(userId, wishlistId, wishlistItemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item moved to cart"})
}

func (wc *WishlistController) GetWishlistItems(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	wishlistId := c.Query("wishlistId")
	if wishlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wishlistId is missing in Query"})
		return
	}

	wishlist, err := wc.WishlistService.GetWishlistItems(userId, wishlistId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wishlist": wishlist})
}
