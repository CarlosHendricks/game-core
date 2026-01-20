package game

import "math"

// CheckBallPaddleCollision checks if the ball collides with a paddle using AABB
func CheckBallPaddleCollision(ball *Ball, paddle *Paddle) bool {
	// Get ball bounds
	ballLeft := ball.X - ball.Radius
	ballRight := ball.X + ball.Radius
	ballTop := ball.Y - ball.Radius
	ballBottom := ball.Y + ball.Radius

	// Get paddle bounds
	paddleLeft := paddle.X
	paddleRight := paddle.X + paddle.Width
	paddleTop := paddle.Y
	paddleBottom := paddle.Y + paddle.Height

	// AABB collision detection
	return ballRight > paddleLeft &&
		ballLeft < paddleRight &&
		ballBottom > paddleTop &&
		ballTop < paddleBottom
}

// HandleBallPaddleCollision handles the ball bouncing off a paddle
func HandleBallPaddleCollision(ball *Ball, paddle *Paddle) {
	// Reverse X direction
	ball.VelocityX = -ball.VelocityX

	// Calculate where on the paddle the ball hit (0 = top, 1 = bottom)
	relativeIntersectY := (paddle.Y + (paddle.Height / 2)) - ball.Y
	normalizedIntersectY := relativeIntersectY / (paddle.Height / 2)

	// Calculate bounce angle (max 60 degrees)
	bounceAngle := normalizedIntersectY * (math.Pi / 3)

	// Calculate new velocity based on bounce angle
	speed := math.Sqrt(ball.VelocityX*ball.VelocityX + ball.VelocityY*ball.VelocityY)
	
	// Determine direction based on which paddle was hit
	direction := 1.0
	if ball.VelocityX < 0 {
		direction = -1.0
	}

	ball.VelocityX = direction * speed * math.Cos(bounceAngle)
	ball.VelocityY = -speed * math.Sin(bounceAngle)

	// Slightly increase speed on each hit (max 1.5x original speed)
	maxSpeed := BallSpeed * 1.5
	if speed < maxSpeed {
		ball.VelocityX *= 1.05
		ball.VelocityY *= 1.05
	}

	// Move ball out of paddle to prevent double collision
	if direction > 0 {
		ball.X = paddle.X + paddle.Width + ball.Radius
	} else {
		ball.X = paddle.X - ball.Radius
	}
}

// UpdateBallPosition updates the ball position and handles wall collisions
func UpdateBallPosition(ball *Ball, fieldHeight float64) {
	ball.X += ball.VelocityX
	ball.Y += ball.VelocityY

	// Top and bottom wall collisions
	if ball.Y-ball.Radius <= 0 {
		ball.Y = ball.Radius
		ball.VelocityY = -ball.VelocityY
	}
	if ball.Y+ball.Radius >= fieldHeight {
		ball.Y = fieldHeight - ball.Radius
		ball.VelocityY = -ball.VelocityY
	}
}

// CheckGoal checks if the ball has gone past the paddles (scoring)
// Returns: 0 = no goal, 1 = player scored, 2 = AI scored
func CheckGoal(ball *Ball, fieldWidth float64) int {
	if ball.X-ball.Radius <= 0 {
		return 2 // AI scored (ball went past player's paddle)
	}
	if ball.X+ball.Radius >= fieldWidth {
		return 1 // Player scored (ball went past AI's paddle)
	}
	return 0
}
