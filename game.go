package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenW = 1280
	screenH = 720
)

const paddleWidth = 15
const paddleHeight = 150
const ballSize = 18
const ballSpeed = 400.0
const winScore = 10

type Paddle struct {
	rect  rl.Rectangle
	speed float32
}

type Ball struct {
	position  rl.Vector2
	direction rl.Vector2
	speed     float32
}

type Sounds struct {
	hit  rl.Sound
	edge rl.Sound
	top  rl.Sound
}

const (
	MAIN_MENU = iota
	PREGAME
	GAME
	GAME_OVER
	EXIT_WINDOW
)

type gameScreen int

type GameState struct {
	leftScore    int
	rightScore   int
	isPaused     bool
	sounds       Sounds
	aiPlayer     bool
	currentScene gameScreen
	prevScene    gameScreen
}

func main() {
	SearchAndSetResourceDir("resources")
	rl.SetTraceLogLevel(rl.LogAll)
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagMsaa4xHint)
	rl.InitWindow(screenW, screenH, "Pong Game")
	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	rl.InitAudioDevice()
	var sn_beep rl.Sound = rl.LoadSound("sounds/ping_pong_8bit_beeep.ogg")
	var sn_peep rl.Sound = rl.LoadSound("sounds/ping_pong_8bit_peeeeeep.ogg")
	var sn_plop rl.Sound = rl.LoadSound("sounds/ping_pong_8bit_plop.ogg")

	state := GameState{
		leftScore:    0,
		rightScore:   0,
		isPaused:     false,
		sounds:       Sounds{top: sn_beep, edge: sn_peep, hit: sn_plop},
		aiPlayer:     true,
		currentScene: MAIN_MENU,
	}

	rl.SetExitKey(rl.KeyNull) // Disable KEY_ESCAPE to close window, X-button still works

	exitWindow := false // Flag to set window to exit

	leftPaddle := Paddle{
		rl.Rectangle{
			X: 50, Y: (screenH / 2) - (paddleHeight / 2),
			Width: paddleWidth, Height: paddleHeight}, 500.0}
	rightPaddle := Paddle{
		rl.Rectangle{X: screenW - 50 - paddleWidth,
			Y:     (screenH / 2) - (paddleHeight / 2),
			Width: paddleWidth, Height: paddleHeight}, 500.0}
	ball := Ball{
		rl.Vector2{X: (screenW / 2), Y: (screenH / 2)},
		rl.Vector2{X: 1.0, Y: 1.0}, ballSpeed,
	}
	ResetBall(&ball)

	for !exitWindow {
		if rl.WindowShouldClose() || rl.IsKeyPressed(rl.KeyEscape) {
			state.prevScene = state.currentScene
			state.currentScene = EXIT_WINDOW
		}

		if state.currentScene == EXIT_WINDOW {
			if rl.IsKeyPressed(rl.KeyY) {
				exitWindow = true
			} else if rl.IsKeyPressed(rl.KeyN) {
				state.currentScene = state.prevScene
			}
		}

		switch state.currentScene {
		case MAIN_MENU:
			if rl.IsKeyPressed(rl.KeyEnter) {
				state.currentScene = PREGAME
				state.aiPlayer = true
			} else if rl.IsKeyPressed(rl.KeySpace) {
				state.currentScene = GAME
				state.aiPlayer = false
			}
		case PREGAME:
			if rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyEnter) ||
				rl.IsKeyPressed(rl.KeyW) || rl.IsKeyPressed(rl.KeyUp) ||
				rl.IsKeyPressed(rl.KeyS) || rl.IsKeyPressed(rl.KeyDown) {
				state.currentScene = GAME
			}
		case GAME:
			if !state.isPaused {
				GameLogic(&leftPaddle, &rightPaddle, &ball, &state)
			}
			if rl.IsKeyPressed(rl.KeyP) {
				state.isPaused = !state.isPaused
			}
		case GAME_OVER:
			ball.speed = 0.0
			if rl.IsKeyPressed(rl.KeyR) {
				state.leftScore = 0
				state.rightScore = 0
				ResetBall(&ball)
				state.currentScene = GAME
			} else if rl.IsKeyPressed(rl.KeyM) {
				state.currentScene = MAIN_MENU
			}
		default:
			break
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)

		switch state.currentScene {
		case EXIT_WINDOW:
			exit_str := "Are you sure you want to exit program? [Y/N]"
			rl.DrawText(exit_str, int32(rl.GetScreenWidth())/2-rl.MeasureText(exit_str, 30)/2, int32(rl.GetScreenHeight())/2, 30, rl.RayWhite)
		case MAIN_MENU:
			defaultFont := rl.GetFontDefault()
			origin := rl.Vector2{X: 0, Y: 0}
			var rotation float32 = 0.0 // No rotation
			fontSize := float32(40.0)  // Normal size
			spacing := float32(2.0)
			rl.DrawTextPro(defaultFont, "PONG!", rl.Vector2{X: float32(rl.GetScreenWidth())/2 - float32(rl.MeasureText("PONG!", int32(fontSize)))/2, Y: float32(rl.GetScreenHeight())/2 - 80}, origin, rotation, fontSize, spacing, rl.RayWhite)
			rl.DrawTextPro(defaultFont, "Play with AI (press Enter)", rl.Vector2{X: float32(rl.GetScreenWidth())/2 - float32(rl.MeasureText("Play with AI (press Enter)", 20))/2, Y: float32(rl.GetScreenHeight()) / 2}, origin, rotation, 20, spacing, rl.RayWhite)
			rl.DrawTextPro(defaultFont, "Local Two Player (press Space)", rl.Vector2{X: float32(rl.GetScreenWidth())/2 - float32(rl.MeasureText("Local Two Player (press Space)", 20))/2, Y: float32(rl.GetScreenHeight())/2 + 40}, origin, rotation, 20, spacing, rl.RayWhite)
			rl.DrawTextPro(defaultFont, "Exit (Press Escape)", rl.Vector2{X: float32(rl.GetScreenWidth())/2 - float32(rl.MeasureText("Exit (Press Escape)", 20))/2, Y: float32(rl.GetScreenHeight())/2 + 80}, origin, rotation, 20, spacing, rl.RayWhite)
		case PREGAME:
			rl.DrawRectangleRec(leftPaddle.rect, rl.LightGray)
			rl.DrawRectangleRec(rightPaddle.rect, rl.LightGray)
			DrawDashedLine(rl.LightGray)
			rl.DrawText(fmt.Sprintf("%d", state.leftScore), screenW/4, 20, 80, rl.LightGray)
			rl.DrawText(fmt.Sprintf("%d", state.rightScore), 3*screenW/4, 20, 80, rl.LightGray)
			rl.DrawText(fmt.Sprintf("Press to start game."), (screenW/2)-(screenW/4), screenH/3, 60, rl.RayWhite)
		case GAME:
			rl.DrawRectangleRec(leftPaddle.rect, rl.RayWhite)
			rl.DrawRectangleRec(rightPaddle.rect, rl.RayWhite)
			rl.DrawCircleV(ball.position, ballSize/2, rl.RayWhite)
			DrawDashedLine(rl.RayWhite)
			rl.DrawText(fmt.Sprintf("%d", state.leftScore), screenW/4, 20, 80, rl.RayWhite)
			rl.DrawText(fmt.Sprintf("%d", state.rightScore), 3*screenW/4, 20, 80, rl.RayWhite)
			if state.isPaused {
				rl.DrawRectangleRec(leftPaddle.rect, rl.LightGray)
				rl.DrawRectangleRec(rightPaddle.rect, rl.LightGray)
				DrawDashedLine(rl.LightGray)
				rl.DrawText(fmt.Sprintf("Paused."), int32(rl.GetScreenWidth())/2-rl.MeasureText("Paused", 60)/2, screenH/2, 60, rl.Red)
			}
		case GAME_OVER:
			rl.DrawRectangleRec(leftPaddle.rect, rl.LightGray)
			rl.DrawRectangleRec(rightPaddle.rect, rl.LightGray)
			DrawDashedLine(rl.LightGray)
			rl.DrawText(fmt.Sprintf("%d", state.leftScore), screenW/4, 20, 80, rl.LightGray)
			rl.DrawText(fmt.Sprintf("%d", state.rightScore), 3*screenW/4, 20, 80, rl.LightGray)
			rl.DrawText(fmt.Sprintf("Press R to restart game."), (screenW/2)-(screenW/4), screenH/2, 60, rl.RayWhite)
		}

		rl.EndDrawing()
	}

	// De-Initialization
	rl.UnloadSound(sn_beep)
	rl.UnloadSound(sn_peep)
	rl.UnloadSound(sn_plop)
	rl.CloseAudioDevice() // Close audio device

	rl.CloseWindow() // Close window and OpenGL context
}

func GameLogic(leftPaddle *Paddle, rightPaddle *Paddle, ball *Ball, state *GameState) {
	if state.aiPlayer {
		if ball.direction.X > 0 {
			// If the ball is moving away from the left paddle
			// Move the paddle slowly towards the center of the screen
			centerY := (screenH - leftPaddle.rect.Height) / 2
			distanceToCenter := (float64)(centerY - (leftPaddle.rect.Y + leftPaddle.rect.Height/2))

			if math.Abs(distanceToCenter) > 10.0 { // Only move if the distance is significant
				if distanceToCenter > 0 {
					leftPaddle.rect.Y += (leftPaddle.speed) * rl.GetFrameTime() // Move down slowly
				} else {
					leftPaddle.rect.Y -= (leftPaddle.speed) * rl.GetFrameTime() // Move up slowly
				}
			}
		} else {

			// Calculate the target position based on the ball's position
			targetY := ball.position.Y

			// Calculate the distance to the target position
			distance := float64(targetY - (leftPaddle.rect.Y + leftPaddle.rect.Height/2))

			// Move the paddle towards the ball's position
			if math.Abs(distance) > 1.0 { // Only move if the distance is significant
				if distance > 0 {
					leftPaddle.rect.Y += leftPaddle.speed * rl.GetFrameTime() // Move down
				} else {
					leftPaddle.rect.Y -= leftPaddle.speed * rl.GetFrameTime() // Move up
				}
			}

			// Clamp the paddle's position to stay within the game boundaries
			if leftPaddle.rect.Y < 0 {
				leftPaddle.rect.Y = 0
			} else if leftPaddle.rect.Y > screenH-leftPaddle.rect.Height {
				leftPaddle.rect.Y = screenH - leftPaddle.rect.Height
			}
		}

		// Extended Player Input
		if (rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp)) && rightPaddle.rect.Y > 0 {
			rightPaddle.rect.Y -= rightPaddle.speed * rl.GetFrameTime()
		}
		if (rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown)) && rightPaddle.rect.Y < screenH-paddleHeight {
			rightPaddle.rect.Y += rightPaddle.speed * rl.GetFrameTime()
		}
	} else {
		// Two PlayerInput
		if rl.IsKeyDown(rl.KeyW) && leftPaddle.rect.Y > 0 {
			leftPaddle.rect.Y -= leftPaddle.speed * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyS) && leftPaddle.rect.Y < screenH-paddleHeight {
			leftPaddle.rect.Y += leftPaddle.speed * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyUp) && rightPaddle.rect.Y > 0 {
			rightPaddle.rect.Y -= rightPaddle.speed * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyDown) && rightPaddle.rect.Y < screenH-paddleHeight {
			rightPaddle.rect.Y += rightPaddle.speed * rl.GetFrameTime()
		}
	}

	// Update ball position
	ball.position.X += ball.direction.X * ball.speed * rl.GetFrameTime()
	ball.position.Y += ball.direction.Y * ball.speed * rl.GetFrameTime()

	// Ball collision with top and bottom
	if ball.position.Y <= 0 || ball.position.Y >= screenH-ballSize {
		ball.direction.Y *= -1 // Reverse Y direction
		rl.PlaySound(state.sounds.top)
	}

	// Ball collision with paddles
	if rl.CheckCollisionCircleRec(ball.position, ballSize/2, leftPaddle.rect) ||
		rl.CheckCollisionCircleRec(ball.position, ballSize/2, rightPaddle.rect) {
		ball.direction.X *= -1 // Reverse X direction
		ball.speed += ballSpeed / 10.0
		rl.PlaySound(state.sounds.hit)
	}

	// Scoring
	if ball.position.X < 0 {
		rl.PlaySound(state.sounds.edge)
		state.rightScore++ // Right player scores
		ResetBall(ball)
	}
	if ball.position.X > screenW {
		rl.PlaySound(state.sounds.edge)
		state.leftScore++ // Left player scores
		ResetBall(ball)
	}

	if state.leftScore >= winScore || state.rightScore >= winScore {
		state.currentScene = GAME_OVER
	}

	// TraceLog(LOG_DEBUG, "After Update: x: %f, y: %f", ball.direction.X, ball.direction.Y)
}

func ResetBall(ball *Ball) {
	ball.position = rl.Vector2{X: screenW / 2, Y: screenH / 2}
	if rl.GetRandomValue(0, 1) == 0 { // Randomize initial direction
		ball.direction.X = 1.0
	} else {
		ball.direction.X = -1.0
	}
	if rl.GetRandomValue(-1, 1) < 0 { // Randomize initial vertical direction
		ball.direction.Y = 1.0
	} else {
		ball.direction.Y = -1.0
	}
	ball.speed = ballSpeed
}

func DrawDashedLine(color rl.Color) {
	start := rl.Vector2{X: screenW / 2, Y: 0}
	end := rl.Vector2{X: screenW / 2, Y: screenH}
	dashLength := float32(10.0)
	gapLength := float32(5.0)
	direction := rl.Vector2Subtract(end, start)
	length := rl.Vector2Length(direction)
	direction = rl.Vector2Normalize(direction)

	var totalLength float32 = 0.0
	for totalLength < length {
		// Calculate the start and end points of the dash
		dashStart := rl.Vector2Add(start, rl.Vector2Scale(direction, totalLength))
		dashEnd := rl.Vector2Add(dashStart, rl.Vector2Scale(direction, dashLength))

		// Draw the dash
		rl.DrawLineEx(dashStart, dashEnd, 5, color)

		// Move the total length by the dash length and gap length
		totalLength += dashLength + gapLength
	}
}
