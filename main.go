package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game constants
const (
	WINDOW_WIDTH   = 1000
	WINDOW_HEIGHT  = 650
	FIRE_RATE     = 0.18
	FALL_RATE     = 1.0
	MAX_BULLETS   = 20
	MAX_ROCKS     = 40
	SCORE_PER_ROCK = 10
	BGM_DURATION  = 75.0
)

func main() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Space Invaders")
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	assets := LoadGameAssets()
	defer assets.UnloadAll()

	rl.SetWindowIcon(*assets.icon)

	sound := LoadGameSounds()
	defer sound.UnloadAllSound()

	// Initialize game state with new function
	gameState := NewGameState()

	InitBullets(MAX_BULLETS, assets.bullet)
	InitRocks(MAX_ROCKS, 100, 4, assets.smallRock)
	ship := NewSpaceShip(rl.Vector2{X: 400, Y: 400})

	var lastFireTime float32 = 0
	var lastFallTime float32 = 0
	var bgmTimer float32 = 0
	
	// Only start playing BGM when game starts
	var bgmStarted bool = false

	for !rl.WindowShouldClose() {
		deltaTime := rl.GetFrameTime()
		
		// Start BGM only when game starts
		if gameState.gamePhase == STATE_PLAYING && !bgmStarted {
			rl.PlaySound(sound.bgm)
			bgmStarted = true
		}
		
		// Update BGM timer only during gameplay
		if gameState.gamePhase == STATE_PLAYING {
			bgmTimer += deltaTime
			if bgmTimer >= BGM_DURATION {
				rl.PlaySound(sound.bgm)
				bgmTimer = 0
			}
		}

		lastFireTime, lastFallTime = updateGame(ship, &sound, &gameState, deltaTime, lastFireTime, lastFallTime)
		drawGame(assets, ship, &gameState)
	}
}
