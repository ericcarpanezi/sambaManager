package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type eventHub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func newEventHub() *eventHub {
	h := &eventHub{clients: map[*websocket.Conn]struct{}{}}
	go h.heartbeat()
	return h
}

func (h *eventHub) Handle(c *gin.Context) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	h.mu.Lock()
	h.clients[conn] = struct{}{}
	h.mu.Unlock()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
			_ = conn.Close()
			return
		}
	}
}

func (h *eventHub) heartbeat() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for now := range ticker.C {
		h.mu.Lock()
		for c := range h.clients {
			if err := c.WriteJSON(gin.H{"event": "heartbeat", "timestamp": now.UTC()}); err != nil {
				_ = c.Close()
				delete(h.clients, c)
			}
		}
		h.mu.Unlock()
	}
}
