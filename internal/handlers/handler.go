package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/haqer0002/avito-shop/internal/middleware"
	"github.com/haqer0002/avito-shop/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Публичные маршруты
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signIn)
	}

	// Защищенные маршруты
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(h.services.Auth))
	{
		user := api.Group("/user")
		{
			user.GET("/info", h.getUserInfo)
			user.POST("/send", h.sendCoins)
		}

		merch := api.Group("/merch")
		{
			merch.POST("/buy/:item", h.buyMerch)
			merch.GET("/list", h.getAllMerch)
		}
	}

	return router
}
