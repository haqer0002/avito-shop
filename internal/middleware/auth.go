package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/haqer0002/avito-shop/internal/service"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
)

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authorizationHeader)
		if header == "" {
			log.Printf("Empty auth header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
			return
		}
		log.Printf("Got auth header: %s", header)

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			log.Printf("Invalid auth header format: %s", header)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}
		log.Printf("Got token: %s", headerParts[1])

		userID, err := authService.ParseToken(headerParts[1])
		if err != nil {
			log.Printf("Error parsing token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Successfully parsed token, got user ID: %d (type: %T)", userID, userID)

		c.Set(userCtx, userID)
		log.Printf("Set user ID in context: %d (type: %T)", userID, userID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (int64, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		log.Printf("User ID not found in context")
		return 0, errors.New("user id not found")
	}
	log.Printf("Got user ID from context: %v (type: %T)", id, id)

	idInt64, ok := id.(int64)
	if !ok {
		log.Printf("User ID is of invalid type: %T, value: %v", id, id)
		return 0, errors.New("user id is of invalid type")
	}
	log.Printf("Successfully converted user ID to int64: %d", idInt64)

	return idInt64, nil
}
