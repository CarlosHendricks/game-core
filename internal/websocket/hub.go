package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// Client represents a connected client
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// Message represents a WebSocket message
type Message struct {
	Type  string      `json:"type"`
	Data  interface{} `json:"data"`
}

var hub *Hub

func init() {
	hub = &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go hub.run()
}

// GetHub returns the singleton hub instance
func GetHub() *Hub {
	return hub
}

// run starts the hub's main loop
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			count := len(h.clients)
			h.mu.Unlock()
			
			log.Printf("Client registered. Total clients: %d", count)
			h.broadcastClientCount()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			count := len(h.clients)
			h.mu.Unlock()
			
			log.Printf("Client unregistered. Total clients: %d", count)
			h.broadcastClientCount()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// broadcastClientCount sends the current client count to all connected clients
func (h *Hub) broadcastClientCount() {
	h.mu.RLock()
	count := len(h.clients)
	h.mu.RUnlock()

	msg := Message{
		Type: "client_count",
		Data: map[string]int{
			"count": count,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling client count: %v", err)
		return
	}

	h.broadcast <- data
}

// ClientCount returns the current number of connected clients
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
