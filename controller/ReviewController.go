package controllers

import (
	"eCommerce/models"
	"eCommerce/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReviewController struct {
	reviewService *services.ReviewService
}

func NewReviewContoller(reviewServices *services.ReviewService) *ReviewController {
	return &ReviewController{
		reviewService: reviewServices,
	}
}

func (rc *ReviewController) SubmitReview(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	productId := c.Query("productId")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing productId in query"})
		return
	}

	var reviewData *models.ReviewData
	if err := c.ShouldBindJSON(&reviewData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := rc.reviewService.SubmitReview(userId, productId, reviewData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review Sumbitted successfully"})
}

func (rc *ReviewController) DeleteReview(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	reviewId := c.Query("reviewId")
	if reviewId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing reviewId in query"})
		return
	}

	err := rc.reviewService.DeleteReview(userId, reviewId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review Deleted successfully"})
}

func (rc *ReviewController) EditReview(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	reviewId := c.Query("reviewId")
	if reviewId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing reviewId in query"})
		return
	}

	var reviewData *models.ReviewData
	if err := c.ShouldBindJSON(&reviewData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := rc.reviewService.EditReview(reviewId, userId, reviewData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review Edited successfully"})
}

func (rc *ReviewController) GetReviews(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	reviews, err := rc.reviewService.GetReviews(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func (rc *ReviewController) GetVendorReviews(c *gin.Context) {
	vendorIdRaw, exists := c.Get("UserId")
	vendorId, ok := vendorIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid vendor ID in context"})
		return
	}

	reviews, err := rc.reviewService.GetVendorReviews(vendorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}
