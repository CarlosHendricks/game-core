package game

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// Client to Server messages
	MsgPlayerInput MessageType = "player_input"
	MsgStartGame   MessageType = "start_game"
	MsgResetGame   MessageType = "reset_game"

	// Server to Client messages
	MsgGameState MessageType = "game_state"
	MsgError     MessageType = "error"
)

// Message represents a generic WebSocket message
type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// InputData represents player input data
type InputData struct {
	Direction float64 `json:"direction"` // -1 (up), 0 (stop), 1 (down)
	PlayerID  int     `json:"playerId,omitempty"` // 1 or 2 (assigned by server)
}

// StateData represents the game state data sent to clients
type StateData struct {
	*GameState
}

// ErrorData represents an error message
type ErrorData struct {
	Message string `json:"message"`
}
