package router

import (
	"aquanote-backend/internal/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) *gin.Engine {
	api := r.Group("/api/v1")
	{
		api.GET("/users", handler.GetUsers)
		api.POST("/users", handler.CreateUser)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "AquaNote"})
	})
	r.GET("/ws", handler.WsHandler)
	r.GET("/sensor/latest", handler.GetLatestHandler)

	return r
}
