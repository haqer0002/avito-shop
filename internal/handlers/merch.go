package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haqer0002/avito-shop/internal/middleware"
)

func (h *Handler) buyMerch(c *gin.Context) {
	merchName := c.Param("item")
	if merchName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "merch name is required"})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.services.Merch.BuyMerch(c.Request.Context(), userID, merchName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) getAllMerch(c *gin.Context) {
	items, err := h.services.Merch.GetAllMerch(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}
