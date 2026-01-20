package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rebec/jueguito/game-core/internal/game"
)

// Hub maintains the set of active clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	game       *game.Game
	running    bool
}

// Client represents a connected client
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	playerID int // 1 or 2, assigned when client connects
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
		game:       game.NewGame(),
		running:    false,
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
			
			// Assign player ID (1 or 2, max 2 players)
			count := len(h.clients)
			if count >= 2 {
				h.mu.Unlock()
				// Reject connection if already 2 players
				log.Println("Client rejected: game is full (2 players)")
				client.conn.Close()
				continue
			}
			
			client.playerID = count + 1
			h.clients[client] = true
			count = len(h.clients)
			
			// Update player count in game
			h.game.SetPlayerCount(count)
			
			// Start game loop when first client connects
			if count == 1 && !h.running {
				h.running = true
				h.game.Start(h.BroadcastToAll)
				log.Println("Starting game loop (first client connected)")
			}
			h.mu.Unlock()
			
			log.Printf("Client registered as Player %d. Total clients: %d", client.playerID, count)
			
			// Send current game state to new client
			h.sendGameStateToClient(client)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				log.Printf("Player %d disconnected", client.playerID)
				delete(h.clients, client)
				close(client.send)
				
				// Update player count
				count := len(h.clients)
				h.game.SetPlayerCount(count)
				
				// Reassign player IDs for remaining clients
				if count > 0 {
					newID := 1
					for c := range h.clients {
						c.playerID = newID
						newID++
					}
				}
			}
			count := len(h.clients)
			h.mu.Unlock()
			
			log.Printf("Client unregistered. Total clients: %d", count)

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

// BroadcastToAll sends a message to all connected clients
func (h *Hub) BroadcastToAll(data []byte) {
	h.broadcast <- data
}

// sendGameStateToClient sends the current game state to a specific client
func (h *Hub) sendGameStateToClient(client *Client) {
	state := h.game.GetState()
	msg := game.Message{
		Type: game.MsgGameState,
		Data: state,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling game state: %v", err)
		return
	}

	select {
	case client.send <- data:
	default:
		log.Printf("Failed to send game state to client")
	}
}

// ProcessMessage processes incoming messages from clients
func (h *Hub) ProcessMessage(client *Client, msgType game.MessageType, msgData json.RawMessage) {
	switch msgType {
	case game.MsgPlayerInput:
		var input game.InputData
		if err := json.Unmarshal(msgData, &input); err != nil {
			log.Printf("Error unmarshaling input: %v", err)
			return
		}
		// Use the client's assigned player ID
		h.game.HandlePlayerInput(client.playerID, input.Direction)

	case game.MsgStartGame:
		h.game.StartGame()
		log.Println("Game started by client")

	case game.MsgResetGame:
		h.game.ResetGame()
		log.Println("Game reset by client")

	default:
		log.Printf("Unknown message type: %s", msgType)
	}
}

// Stop stops the hub and game loop
func (h *Hub) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if h.running {
		h.game.Stop()
		h.running = false
	}
}

// ClientCount returns the current number of connected clients
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
