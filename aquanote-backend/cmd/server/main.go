package main

import (
	"aquanote-backend/internal/handler"
	"aquanote-backend/internal/router"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 啟動 MQTT Broker
	handler.StartMQTTBroker()

	// 設定 Gin Web Server
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		// AllowOrigins:    []string{"http://localhost:9001"},
		AllowMethods:    []string{"*"},
		AllowWebSockets: true,
	}))

	router.Setup(r)

	if err := r.Run(":9000"); err != nil {
		log.Fatalf("[HTTP Server] Run error: %v", err)
	}
}
