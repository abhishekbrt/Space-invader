package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

// Add difficulty constants
const (
	DIFFICULTY_INCREASE_INTERVAL = 30.0 // Seconds between difficulty increases
	INITIAL_FALL_RATE            = 1.0
	MIN_FALL_RATE                = 0.2 // Fastest spawn rate
	FALL_RATE_DECREASE           = 0.2 // How much to decrease spawn time each level
)

// Add new game state
const (
	STATE_INTRO = iota
	STATE_PLAYING
	STATE_GAMEOVER
)

// Game assets
type GameAssets struct {
	background rl.Texture2D // Main background image
	starsA     rl.Texture2D // First star layer
	starsB     rl.Texture2D // Second star layer 
	spaceship  rl.Texture2D // Player's ship texture
	bullet     rl.Texture2D // Projectile texture
	smallRock  rl.Texture2D // Asteroid texture
	icon       *rl.Image     // Game window icon
}

type GameSounds struct {
	blast       rl.Sound // bullet shooting sound
	bgm         rl.Sound //background music
	explosion   rl.Sound //asteroid explosion sound
	shipExplode rl.Sound //ship explosion sound
}

type GameState struct {
	score           int32
	highScore       int32
	isGameOver      bool
	difficultyTimer float32
	currentFallRate float32
	difficultyLevel int32
	gamePhase       int32 // Current game phase (intro, playing, game over)
}

func NewGameState() GameState {
	return GameState{
		score:           0,
		highScore:       0,
		isGameOver:      false,
		difficultyTimer: 0,
		currentFallRate: INITIAL_FALL_RATE,
		difficultyLevel: 1,
		gamePhase:       STATE_INTRO,
	}
}

// LoadGameAssets loads and resizes all game textures
func LoadGameAssets() GameAssets {
	assets := GameAssets{}

	// Load textures from embedded files
	assets.background = LoadTextureFromEmbedded("assets/bg.png")
	assets.starsA = LoadTextureFromEmbedded("assets/Stars-A.png")
	assets.starsB = LoadTextureFromEmbedded("assets/Stars-B.png")

	// Load and resize spaceship
	shipImg := rl.LoadImageFromTexture(LoadTextureFromEmbedded("assets/SpaceShip.png"))
	rl.ImageResize(shipImg, 60, 60)
	assets.spaceship = rl.LoadTextureFromImage(shipImg)
	rl.UnloadImage(shipImg)

	// Load and resize bullet
	bulletImg := rl.LoadImageFromTexture(LoadTextureFromEmbedded("assets/plasm.png"))
	rl.ImageResize(bulletImg, 20, 50)
	assets.bullet = rl.LoadTextureFromImage(bulletImg)
	rl.UnloadImage(bulletImg)

	// Load and resize rocks
	smallRockImg := rl.LoadImageFromTexture(LoadTextureFromEmbedded("assets/small-A.png"))
	rl.ImageResize(smallRockImg, 60, 60)
	assets.smallRock = rl.LoadTextureFromImage(smallRockImg)
	rl.UnloadImage(smallRockImg)

	assets.icon=LoadImageFromEmbedded("assets/SpaceShip.png")

	return assets
}

func LoadGameSounds() GameSounds {
	return GameSounds{
		blast:       LoadSoundFromEmbedded("assets/sound/blaster-shot.mp3"),
		bgm:         LoadSoundFromEmbedded("assets/sound/drone-bg.mp3"),
		explosion:   LoadSoundFromEmbedded("assets/sound/explosion2.wav"),
		shipExplode: LoadSoundFromEmbedded("assets/sound/explosion.wav"),
	}
}

func (assets *GameAssets) UnloadAll() {
	rl.UnloadTexture(assets.background)
	rl.UnloadTexture(assets.starsA)
	rl.UnloadTexture(assets.starsB)
	rl.UnloadTexture(assets.spaceship)
	rl.UnloadTexture(assets.bullet)
	rl.UnloadTexture(assets.smallRock)

}

func (s *GameSounds) UnloadAllSound() {
	rl.UnloadSound(s.blast)
	rl.UnloadSound(s.bgm)
	rl.UnloadSound(s.explosion)
	rl.UnloadSound(s.shipExplode)
}

// Add new function to draw intro screen
func drawIntroScreen() {
	textColor := rl.White
	titlePos := rl.Vector2{X: float32(WINDOW_WIDTH)/2 - 250, Y: 150}

	rl.DrawText("SPACE INVADERS", int32(titlePos.X), int32(titlePos.Y), 50, rl.Red)

	// Draw game rules
	rulesStartY := int32(titlePos.Y) + 100
	rl.DrawText("Game Controls:", 350, rulesStartY, 30, textColor)
	rl.DrawText("Arrow keys  : Move Ship", 350, rulesStartY+50, 25, textColor)
	rl.DrawText("SPACE   : Shoot", 350, rulesStartY+90, 25, textColor)
	rl.DrawText("Survive and destroy asteroids!", 300, rulesStartY+150, 25, textColor)

	// Draw start prompt
	rl.DrawText("Press ENTER to Start", int32(WINDOW_WIDTH)/2-200, 500, 30, rl.Green)
}

func updateGame(ship *Spaceship, sound *GameSounds, gamestate *GameState, deltaTime, lastFireTime, lastFallTime float32) (float32, float32) {
	// Handle different game phases
	switch gamestate.gamePhase {
	case STATE_INTRO:
		if rl.IsKeyPressed(rl.KeyEnter) {
			gamestate.gamePhase = STATE_PLAYING
			return 0, 0
		}
		return lastFireTime, lastFallTime

	case STATE_PLAYING:
		if gamestate.isGameOver {
			gamestate.gamePhase = STATE_GAMEOVER
			return lastFireTime, lastFallTime
		}
		// ...existing playing state code...
		if gamestate.isGameOver {
			if rl.IsKeyPressed(rl.KeyR) {
				gamestate.isGameOver = false
				gamestate.score = 0
				gamestate.difficultyTimer = 0
				gamestate.currentFallRate = INITIAL_FALL_RATE
				gamestate.difficultyLevel = 1
				ship.lives = 1 // Reset to 1 life instead of 3
				ship.isAlive = true
				for _, rock := range Rocks {
					rock.projectile.isActive = false
				}
				return 0, 0
			}
			return lastFireTime, lastFallTime
		}

		// Update difficulty
		gamestate.difficultyTimer += deltaTime
		if gamestate.difficultyTimer >= DIFFICULTY_INCREASE_INTERVAL {
			gamestate.difficultyTimer = 0
			gamestate.difficultyLevel++
			gamestate.currentFallRate = float32(math.Max(
				float64(MIN_FALL_RATE),
				float64(INITIAL_FALL_RATE-FALL_RATE_DECREASE*float32(gamestate.difficultyLevel-1)),
			))
		}

		ship.Update(deltaTime, gamestate)

		if ship.isAlive {
			if ship.CheckCollisionWithRocks() {
				rl.PlaySound(sound.shipExplode)
				ship.Die(gamestate)
			}

			if rl.IsKeyDown(rl.KeySpace) && lastFireTime >= FIRE_RATE {
				rl.PlaySound(sound.blast)
				SpawnBullet(ship.positon)
				lastFireTime = 0
			}
		}

		lastFireTime += deltaTime
		lastFallTime += deltaTime

		// Use currentFallRate instead of FALL_RATE
		if lastFallTime >= gamestate.currentFallRate {
			SpawnRock()
			lastFallTime = 0
		}

		UpdateRocks(deltaTime)
		UpdateBullets(deltaTime, gamestate, sound)
		ship.ShipControl()

		return lastFireTime, lastFallTime

	case STATE_GAMEOVER:
		if rl.IsKeyPressed(rl.KeyR) {
			// Reset game
			gamestate.isGameOver = false
			gamestate.score = 0
			gamestate.difficultyTimer = 0
			gamestate.currentFallRate = INITIAL_FALL_RATE
			gamestate.difficultyLevel = 1
			gamestate.gamePhase = STATE_PLAYING
			ship.lives = 1
			ship.isAlive = true
			for _, rock := range Rocks {
				rock.projectile.isActive = false
			}
			return 0, 0
		}
	}
	return lastFireTime, lastFallTime
}

func drawGame(assets GameAssets, ship *Spaceship, state *GameState) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	// Draw background for all states
	rl.DrawTexture(assets.background, 0, 0, rl.White)
	rl.DrawTexture(assets.starsA, 0, 0, rl.White)
	rl.DrawTexture(assets.starsB, 0, 0, rl.White)

	switch state.gamePhase {
	case STATE_INTRO:
		drawIntroScreen()

	case STATE_PLAYING:
		// Draw background layers

		DrawBullets()
		DrawRocks()
		if ship.isAlive {
			rl.DrawTextureV(assets.spaceship, ship.positon, rl.White)
		}

		// Draw UI
		scoreText := fmt.Sprintf("SCORE: %d", state.score)
		rl.DrawText(scoreText, 10, 10, 20, rl.White)

		highScoreText := fmt.Sprintf("HIGH SCORE: %d", state.highScore)
		rl.DrawText(highScoreText, WINDOW_WIDTH-200, 10, 20, rl.White)

		livesText := fmt.Sprintf("LIVES: %d", ship.lives)
		rl.DrawText(livesText, 10, 40, 20, rl.White)

		// Add difficulty level display
		levelText := fmt.Sprintf("LEVEL: %d", state.difficultyLevel)
		rl.DrawText(levelText, 10, 70, 20, rl.White)

		if state.isGameOver {
			textPos := rl.Vector2{
				X: float32(WINDOW_WIDTH)/2 - 100,
				Y: float32(WINDOW_HEIGHT) / 2,
			}
			rl.DrawText("GAME OVER", int32(textPos.X), int32(textPos.Y), 40, rl.Red)
			rl.DrawText("Press R to Restart", int32(textPos.X), int32(textPos.Y)+50, 20, rl.White)
		}

	case STATE_GAMEOVER:
		// ...existing game over drawing code...
		if state.isGameOver {
			textPos := rl.Vector2{
				X: float32(WINDOW_WIDTH)/2 - 100,
				Y: float32(WINDOW_HEIGHT) / 2,
			}
			rl.DrawText("GAME OVER", int32(textPos.X), int32(textPos.Y), 40, rl.Red)
			rl.DrawText("Press R to Restart", int32(textPos.X), int32(textPos.Y)+50, 20, rl.White)
		}
	}

	rl.EndDrawing()
}
