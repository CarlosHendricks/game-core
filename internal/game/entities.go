package game

import "math"

// Paddle represents a player or AI paddle
type Paddle struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Speed  float64 `json:"-"` // Don't send speed to client
}

// Ball represents the game ball
type Ball struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	VelocityX float64 `json:"vx"`
	VelocityY float64 `json:"vy"`
	Radius    float64 `json:"radius"`
	Speed     float64 `json:"-"` // Base speed
}

// GameState represents the complete state of the game
type GameState struct {
	Player1Paddle *Paddle `json:"player1"`
	Player2Paddle *Paddle `json:"player2"`
	Ball          *Ball   `json:"ball"`
	Player1Score  int     `json:"player1Score"`
	Player2Score  int     `json:"player2Score"`
	State         string  `json:"state"` // "waiting", "playing", "gameover"
	Winner        string  `json:"winner,omitempty"` // "player1", "player2", or empty
	FieldWidth    float64 `json:"fieldWidth"`
	FieldHeight   float64 `json:"fieldHeight"`
	PlayerCount   int     `json:"playerCount"` // Number of connected players
}

const (
	// Field dimensions
	FieldWidth  = 800
	FieldHeight = 600

	// Paddle dimensions
	PaddleWidth  = 10
	PaddleHeight = 100
	PaddleSpeed  = 5.0

	// Ball dimensions
	BallRadius = 8
	BallSpeed  = 5.0

	// Game settings
	WinningScore = 5
	PaddleOffset = 20
)

// NewGameState creates a new game state with initial values
func NewGameState() *GameState {
	gs := &GameState{
		Player1Paddle: &Paddle{
			X:      PaddleOffset,
			Y:      FieldHeight/2 - PaddleHeight/2,
			Width:  PaddleWidth,
			Height: PaddleHeight,
			Speed:  PaddleSpeed,
		},
		Player2Paddle: &Paddle{
			X:      FieldWidth - PaddleOffset - PaddleWidth,
			Y:      FieldHeight/2 - PaddleHeight/2,
			Width:  PaddleWidth,
			Height: PaddleHeight,
			Speed:  PaddleSpeed,
		},
		Ball: &Ball{
			X:      FieldWidth / 2,
			Y:      FieldHeight / 2,
			Radius: BallRadius,
			Speed:  BallSpeed,
		},
		Player1Score: 0,
		Player2Score: 0,
		State:        "waiting",
		FieldWidth:   FieldWidth,
		FieldHeight:  FieldHeight,
		PlayerCount:  0,
	}

	gs.ResetBall()
	return gs
}

// ResetBall resets the ball to the center with a random direction
func (gs *GameState) ResetBall() {
	gs.Ball.X = FieldWidth / 2
	gs.Ball.Y = FieldHeight / 2

	// Random angle between -45 and 45 degrees (in radians)
	angle := (math.Pi / 4) * (2*math.Floor(math.Mod(float64(gs.Player1Score+gs.Player2Score), 2)) - 1)
	
	// Alternate direction
	direction := 1.0
	if (gs.Player1Score+gs.Player2Score)%2 == 0 {
		direction = -1.0
	}

	gs.Ball.VelocityX = math.Cos(angle) * gs.Ball.Speed * direction
	gs.Ball.VelocityY = math.Sin(angle) * gs.Ball.Speed
}

// MovePaddle moves a paddle by a direction (-1, 0, 1) with bounds checking
func (p *Paddle) MovePaddle(direction float64, fieldHeight float64) {
	p.Y += direction * p.Speed

	// Bounds checking
	if p.Y < 0 {
		p.Y = 0
	}
	if p.Y+p.Height > fieldHeight {
		p.Y = fieldHeight - p.Height
	}
}
