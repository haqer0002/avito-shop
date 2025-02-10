package handlers

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/haqer0002/avito-shop/internal/middleware"
	"github.com/haqer0002/avito-shop/internal/models"
)

func (h *Handler) getUserInfo(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		log.Printf("Error getting user ID from context: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Getting info for user ID: %d (type: %T)", userID, userID)
	info, err := h.services.User.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Successfully got info for user ID %d: %+v", userID, info)
	c.JSON(http.StatusOK, info)
}

func (h *Handler) sendCoins(c *gin.Context) {
	var input models.SendCoinRequest

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.services.User.SendCoins(c.Request.Context(), userID, input.ToUser, input.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
