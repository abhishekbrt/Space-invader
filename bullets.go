package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bullet struct {
	projectile *Projectile
}

var Bullets []*Bullet

// type BulletPool struct {
// 	pool sync.Pool
// }
// func NewBulletPool(noOfBullets int) *BulletPool {
// 	return &BulletPool{
// 		pool: sync.Pool{
// 			New: func() interface{} {
// 				bullets := make([]*Bullets, noOfBullets)
// 				for i := 0; i < noOfBullets; i++ {
// 					bullets[i] = &Bullets{}
// 				}
// 				return bullets
// 			},
// 		},
// 	}
// }
// func NewBulletPool() *BulletPool {
// 	return &BulletPool{
// 		pool: sync.Pool{
// 			New: func() interface{} {
// 				return &Bullets{}
// 			},
// 		},
// 	}
// }

func InitBullets(maxBullet int, texture rl.Texture2D) {
	Bullets = make([]*Bullet, maxBullet)
	for i := 0; i < maxBullet; i++ {
		Bullets[i] = &Bullet{
			projectile: &Projectile{
				isActive: false,
				position: rl.Vector2{},
				speed:    400, //pixels per second
				texture:  texture,
			},
		}
	}
}

func SpawnBullet(startPos rl.Vector2) {
	for _, bullet := range Bullets {
		if !bullet.projectile.isActive {
			bullet.projectile.isActive = true
			bullet.projectile.position = startPos
			break
		}
	}
}

// Add GameState parameter to UpdateBullets
func UpdateBullets(deltatime float32, gameState *GameState,sound *GameSounds) {
	for _, bullet := range Bullets {
		if bullet.projectile.isActive {
			// Move bullet
			bullet.projectile.position.Y -= bullet.projectile.speed * deltatime

			// Check collision with rocks
			// bulletCenter, bulletRadius := bullet.projectile.GetBoundingCircle(rl.Vector2{X: 60, Y: 60})

			// Check collision with each active rock
			for _, rock := range Rocks {
				if rock.projectile.isActive {
					rockCenter, rockRadius := rock.projectile.GetBoundingCircle(rl.Vector2{X: 55, Y: 55})

					//bullet size {width:20,height:50}
					if CheckCircleRectangleCollision(rockCenter, rockRadius, bullet.projectile.position, rl.Vector2{X: 15, Y: 50}) {
						// Deactivate bullet
						bullet.projectile.isActive = false
						rl.PlaySound(sound.explosion)

						
						rock.lifeRemaining--
						if rock.lifeRemaining <= 0 {
							// Update score when rock is destroyed
							gameState.score += SCORE_PER_ROCK
							if gameState.score > gameState.highScore {
								gameState.highScore = gameState.score
							}
							rock.projectile.isActive = false
							rock.projectile.position.Y = -50
						} else if rock.lifeRemaining == 1 {
							rock.color = rl.Red
						}
						break
					}
				}
			}

			// Deactivate if off-screen
			if bullet.projectile.position.Y < 0 {
				bullet.projectile.isActive = false
			}
		}
	}
}

func DrawBullets() {
	for _, bullet := range Bullets {
		if bullet.projectile.isActive {
			position := rl.Vector2{
				X: bullet.projectile.position.X + 20,
				Y: bullet.projectile.position.Y,
			}
			rl.DrawTextureV(bullet.projectile.texture, position, rl.White)
		}
	}
}
