package controllers

import (
	"eCommerce/models"
	"eCommerce/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService,
	}
}

func (uc *UserController) Login(c *gin.Context) {
	var credentials models.Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, accessToken, refreshToken, err := uc.UserService.Login(&credentials)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie(
		"refreshToken", // cookie name
		refreshToken,   // value
		7*24*60*60,     // maxAge in seconds (7 days)
		"/",            // path
		"",             // domain (empty = current domain)
		false,          // secure (set true in production with HTTPS)
		true,           // httpOnly (can't be accessed by JS)
	)

	c.JSON(http.StatusOK, gin.H{
		"user": models.SanitizedUser{
			UserId:      user.UserId,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Verified:    user.Verified,
		},
		"token": accessToken,
	})
}

func (uc *UserController) Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := uc.UserService.Signup(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent. Please check your inbox."})
}

func (uc *UserController) VerifyUser(c *gin.Context) {
	verificationToken := c.Query("verificationToken")
	if verificationToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "verificationToken is needed"})
		return
	}

	err := uc.UserService.VerifyUser(verificationToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Verified"})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newUser, accessToken, refreshToken, err := uc.UserService.UpdateUser(&user, userId)
	if err != nil {
		status := http.StatusBadRequest
		switch {
		case strings.Contains(err.Error(), "failed to"):
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if accessToken == "" && refreshToken == "" {
		c.JSON(http.StatusOK, gin.H{"user": models.SanitizedUser{
			UserId:      newUser.UserId,
			Name:        newUser.Name,
			Email:       newUser.Email,
			PhoneNumber: newUser.PhoneNumber,
			Verified:    newUser.Verified,
		}})
		return
	}

	c.SetCookie(
		"refreshToken", // cookie name
		refreshToken,   // value
		7*24*60*60,     // maxAge in seconds (7 days)
		"/",            // path
		"",             // domain (empty = current domain)
		false,          // secure (set true in production with HTTPS)
		true,           // httpOnly (can't be accessed by JS)
	)
	c.JSON(http.StatusOK, gin.H{
		"user": models.SanitizedUser{
			UserId:      newUser.UserId,
			Name:        newUser.Name,
			Email:       newUser.Email,
			PhoneNumber: newUser.PhoneNumber,
			Verified:    newUser.Verified,
		},
		"token": accessToken,
	})
}

func (uc *UserController) GetUserProfile(c *gin.Context) {
	userIdRaw, exists := c.Get("UserId")
	userId, ok := userIdRaw.(string)
	if !exists || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	user, err := uc.UserService.GetUserProfile(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": models.SanitizedUser{
			UserId:      user.UserId,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Verified:    user.Verified,
		},
	})
}
