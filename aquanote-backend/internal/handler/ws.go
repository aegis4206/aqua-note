package handler

import (
	models "aquanote-backend/internal/model"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex sync.Mutex
)

var (
	latestData  *models.SensorData
	latestMutex sync.RWMutex
)

func Broadcast(payload []byte) {
	// 解析MQTT Payload 並更新全域最新數據
	var data models.SensorData
	if err := json.Unmarshal(payload, &data); err == nil {
		latestMutex.Lock()
		latestData = &data
		latestMutex.Unlock()
	}
	log.Printf("[MQTT] Received data: %s", string(payload))

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			log.Printf("[WS] Write error: %v", err)
			clientsMutex.Lock()
			conn.Close()
			delete(clients, conn)
			clientsMutex.Unlock()
		}
	}
}

func GetLatestHandler(c *gin.Context) {
	latestMutex.RLock()
	defer latestMutex.RUnlock()

	if latestData == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "No sensor data received yet",
		})
		return
	}

	c.JSON(http.StatusOK, latestData)
}

func WsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] Upgrade error: %v", err)
		return
	}
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()
	log.Printf("[WS] Client connected. Total: %d", len(clients))

	// latestMutex.RLock()
	// if latestData != nil {
	// 	if snapshot, err := json.Marshal(latestData); err == nil {
	// 		conn.WriteMessage(websocket.TextMessage, snapshot)
	// 	}
	// }
	// latestMutex.RUnlock()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
			conn.Close()
			log.Printf("[WS] Client disconnected. Total: %d", len(clients))
			break
		}
	}
}
