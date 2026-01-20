package game

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// Game represents the game instance
type Game struct {
	State         *GameState
	mu            sync.RWMutex
	running       bool
	tickRate      time.Duration
	lastUpdate    time.Time
	player1Input  float64 // Player 1 input direction
	player2Input  float64 // Player 2 input direction
}

const (
	TicksPerSecond = 60
	StateUpdateRate = 20 // Send state 20 times per second
)

// NewGame creates a new game instance
func NewGame() *Game {
	return &Game{
		State:      NewGameState(),
		tickRate:   time.Second / TicksPerSecond,
		lastUpdate: time.Now(),
	}
}

// Start starts the game loop
func (g *Game) Start(broadcastFunc func([]byte)) {
	g.mu.Lock()
	if g.running {
		g.mu.Unlock()
		return
	}
	g.running = true
	g.mu.Unlock()

	log.Println("Game loop started")
	go g.gameLoop(broadcastFunc)
}

// Stop stops the game loop
func (g *Game) Stop() {
	g.mu.Lock()
	g.running = false
	g.mu.Unlock()
	log.Println("Game loop stopped")
}

// gameLoop is the main game loop running at 60 TPS
func (g *Game) gameLoop(broadcastFunc func([]byte)) {
	ticker := time.NewTicker(g.tickRate)
	defer ticker.Stop()

	stateUpdateCounter := 0
	stateUpdateInterval := TicksPerSecond / StateUpdateRate // Send state every N ticks

	for {
		<-ticker.C

		g.mu.Lock()
		if !g.running {
			g.mu.Unlock()
			return
		}

		// Update game state
		g.update()

		// Send state update at reduced rate
		stateUpdateCounter++
		if stateUpdateCounter >= stateUpdateInterval {
			stateUpdateCounter = 0
			g.broadcastState(broadcastFunc)
		}

		g.mu.Unlock()
	}
}

// update updates the game state for one tick
func (g *Game) update() {
	if g.State.State != "playing" {
		return
	}

	// Update player 1 paddle based on input
	if g.player1Input != 0 {
		g.State.Player1Paddle.MovePaddle(g.player1Input, g.State.FieldHeight)
	}

	// Update player 2 paddle based on input
	if g.player2Input != 0 {
		g.State.Player2Paddle.MovePaddle(g.player2Input, g.State.FieldHeight)
	}

	// Update ball position
	UpdateBallPosition(g.State.Ball, g.State.FieldHeight)

	// Check paddle collisions
	if CheckBallPaddleCollision(g.State.Ball, g.State.Player1Paddle) {
		HandleBallPaddleCollision(g.State.Ball, g.State.Player1Paddle)
	}
	if CheckBallPaddleCollision(g.State.Ball, g.State.Player2Paddle) {
		HandleBallPaddleCollision(g.State.Ball, g.State.Player2Paddle)
	}

	// Check for goals
	goal := CheckGoal(g.State.Ball, g.State.FieldWidth)
	if goal != 0 {
		if goal == 1 {
			g.State.Player1Score++
			log.Printf("Player 1 scored! Score: %d - %d", g.State.Player1Score, g.State.Player2Score)
		} else {
			g.State.Player2Score++
			log.Printf("Player 2 scored! Score: %d - %d", g.State.Player1Score, g.State.Player2Score)
		}

		// Check for game over
		if g.State.Player1Score >= WinningScore {
			g.State.State = "gameover"
			g.State.Winner = "player1"
			log.Println("Game Over: Player 1 wins!")
		} else if g.State.Player2Score >= WinningScore {
			g.State.State = "gameover"
			g.State.Winner = "player2"
			log.Println("Game Over: Player 2 wins!")
		} else {
			// Reset ball for next round
			g.State.ResetBall()
		}
	}
}

// broadcastState sends the current game state to all clients
func (g *Game) broadcastState(broadcastFunc func([]byte)) {
	msg := Message{
		Type: MsgGameState,
		Data: g.State,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling game state: %v", err)
		return
	}

	broadcastFunc(data)
}

// HandlePlayerInput handles player input messages
func (g *Game) HandlePlayerInput(playerID int, direction float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Clamp direction to -1, 0, or 1
	clampedDirection := 0.0
	if direction < -0.5 {
		clampedDirection = -1
	} else if direction > 0.5 {
		clampedDirection = 1
	}

	// Update the appropriate player's input
	if playerID == 1 {
		g.player1Input = clampedDirection
	} else if playerID == 2 {
		g.player2Input = clampedDirection
	}
}

// StartGame starts a new game
func (g *Game) StartGame() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.State.State == "playing" {
		return
	}

	// Only start if we have 2 players
	if g.State.PlayerCount < 2 {
		log.Println("Cannot start game: need 2 players")
		return
	}

	log.Println("Starting new game")
	g.State.State = "playing"
	g.State.Player1Score = 0
	g.State.Player2Score = 0
	g.State.Winner = ""
	g.State.ResetBall()
}

// ResetGame resets the game to initial state
func (g *Game) ResetGame() {
	g.mu.Lock()
	defer g.mu.Unlock()

	log.Println("Resetting game")
	playerCount := g.State.PlayerCount // Preserve player count
	g.State = NewGameState()
	g.State.PlayerCount = playerCount
	g.player1Input = 0
	g.player2Input = 0
}

// SetPlayerCount updates the number of connected players
func (g *Game) SetPlayerCount(count int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.State.PlayerCount = count
}

// GetState returns a copy of the current game state
func (g *Game) GetState() *GameState {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	stateCopy := *g.State
	return &stateCopy
}
