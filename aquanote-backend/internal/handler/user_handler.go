package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get users",
	})
}

func CreateUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
	})
}
