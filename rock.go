package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ROCK_SIZE_X = 60
	ROCK_SIZE_Y = 60
)

// SmallRock represents a small asteroid in the game
type Rock struct {
	projectile    *Projectile
	lifeRemaining uint8    // Remaining hits before destruction
	color         rl.Color // Add color property
	colorTimer    float32  // Timer for color transition
}

var Rocks []*Rock

func InitRocks(maxRock int, speed float32, lifeRemaining uint8, texture rl.Texture2D) {
	Rocks = make([]*Rock, maxRock)
	for i := 0; i < maxRock; i++ {
		Rocks[i] = &Rock{
			projectile: &Projectile{
				isActive: false,
				position: rl.Vector2{X: 0, Y: -50}, // Start above screen
				speed:    speed,
				texture:  texture,
			},
			lifeRemaining: lifeRemaining,
			color:         rl.White, // Initial color
			colorTimer:    0,
		}
	}
}

func SpawnRock() {
	xPos := float32(rl.GetRandomValue(100, 900))
	for _, rock := range Rocks {
		if !rock.projectile.isActive {
			rock.projectile.isActive = true
			rock.projectile.position = rl.Vector2{X: xPos, Y: -50}
			rock.lifeRemaining = 2 // Reset life when spawning new rock
			rock.color = rl.White  // Reset color when spawning
			rock.colorTimer = 0
			break
		}
	}
}

func UpdateRocks(deltatime float32) {
	activeRockCount := 0
	for _, rock := range Rocks {
		if rock.projectile.isActive {
			activeRockCount++
			rock.projectile.position.Y += deltatime * rock.projectile.speed

		}

		// deactivate if it goes offscreen
		if rock.projectile.position.Y > 650 || rock.lifeRemaining <= 0 {
			rock.projectile.isActive = false
			rock.projectile.position.Y = -50
		}
	}

	// Ensure there's always at least one rock active
	if activeRockCount == 0 {
		SpawnRock()
	}
}

func DrawRocks() {
	for _, rock := range Rocks {
		if rock.projectile.isActive {
			rl.DrawTextureV(rock.projectile.texture, rock.projectile.position, rock.color)
		}
	}
}

