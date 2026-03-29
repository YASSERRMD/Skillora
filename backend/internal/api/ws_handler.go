package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/skillora/backend/internal/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Simple CORS for dev
	},
}

// WSHandler upgrades a specific HTTP connection to a WebSocket.
func WSHandler(h *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		h.Register(userID, conn)
		
		// Keep connection alive until error or disconnect.
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				h.Unregister(userID)
				break
			}
		}
	}
}
