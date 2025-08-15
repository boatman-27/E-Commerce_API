package middlewares

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(c *gin.Context) (claims jwt.MapClaims, err error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
		c.Abort()
		return nil, err
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header format must be Bearer {token}"})
		c.Abort()
		return nil, err
	}

	tokenStr := parts[1]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid or malformed token",
			"error":   err.Error(),
		})
		c.Abort()
		return nil, err
	}

	var ok bool
	claims, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
		c.Abort()
		return nil, err
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid expiration in token"})
		c.Abort()
		return nil, err
	}

	if float64(time.Now().Unix()) > exp {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Access Token has expired"})
		c.Abort()
		return nil, err
	}

	return claims, nil
}

func RequireAuth(c *gin.Context) {
	claims, err := Auth(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Authorization Error",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}

	email, emailOk := claims["Email"].(string)
	userId, userIdOk := claims["UserId"].(string)

	if !emailOk || !userIdOk {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing or invalid token claims"})
		c.Abort()
		return
	}

	c.Set("Email", email)
	c.Set("UserId", userId)

	c.Next()
}

func VendorAuth(c *gin.Context) {
	claims, err := Auth(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Authorization Error",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}

	role, ok := claims["Role"].(string)
	if !ok || role != "vendor" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Admin access required"})
		c.Abort()
		return
	}

	email, emailOk := claims["Email"].(string)
	userId, userIdOk := claims["UserId"].(string)

	if !emailOk || !userIdOk {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing or invalid token claims"})
		c.Abort()
		return
	}

	c.Set("Email", email)
	c.Set("UserId", userId)

	c.Next()
}
