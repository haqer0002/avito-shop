package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haqer0002/avito-shop/internal/models"
)

func (h *Handler) signIn(c *gin.Context) {
	var input models.AuthRequest

	if err := c.BindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
		return
	}

	log.Printf("Received sign-up request for user: %s", input.Username)

	// Проверяем существование пользователя
	_, err := h.services.Auth.GetUserByUsername(c.Request.Context(), input.Username)
	if err != nil {
		log.Printf("User not found, creating new user: %s", input.Username)
		// Пользователь не найден, создаем нового
		if err := h.services.Auth.CreateUser(c.Request.Context(), input.Username, input.Password); err != nil {
			log.Printf("Error creating user: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Successfully created user: %s", input.Username)
	} else {
		log.Printf("User already exists: %s", input.Username)
	}

	// Генерируем токен (для нового или существующего пользователя)
	token, err := h.services.Auth.GenerateToken(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	log.Printf("Successfully generated token for user: %s", input.Username)
	c.JSON(http.StatusOK, models.AuthResponse{Token: token})
}
