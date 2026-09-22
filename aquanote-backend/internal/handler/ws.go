package handler

import (
	model "aquanote-backend/internal/model"
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
	clientsMutex sync.RWMutex
)

var (
	latestData  *model.TemperatureLog
	latestMutex sync.RWMutex
)

func Broadcast(data model.TemperatureLog) {
	latestMutex.Lock()
	latestData = &data
	latestMutex.Unlock()

	clientsMutex.RLock()
	if len(clients) == 0 {
		clientsMutex.RUnlock()
		return
	}
	targets := make([]*websocket.Conn, 0, len(clients))
	for conn := range clients {
		targets = append(targets, conn)
	}
	clientsMutex.RUnlock()

	var deadConns []*websocket.Conn
	for _, conn := range targets {
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("[WS] Write error: %v", err)
			conn.Close()
			deadConns = append(deadConns, conn)
		}
	}

	if len(deadConns) > 0 {
		clientsMutex.Lock()
		for _, conn := range deadConns {
			delete(clients, conn)
		}
		clientsMutex.Unlock()
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
