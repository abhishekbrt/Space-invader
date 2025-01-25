package main

import rl "github.com/gen2brain/raylib-go/raylib"

// Movement constants
const (
    SHIP_SPEED    = 5.0
    MIN_X         = 30
    MAX_X         = 920
    MIN_Y         = 30
    MAX_Y         = 570
)

// Spaceship represents the player's vessel
type Spaceship struct {
    positon rl.Vector2
    lives    int32
    isAlive  bool
    respawnTimer float32
}

// NewSpaceShip creates a new spaceship at the given position
func NewSpaceShip(positon rl.Vector2) *Spaceship {
    return &Spaceship{
        positon: positon,
        lives: 1,        // Changed from 3 to 1
        isAlive: true,
        respawnTimer: 0,
    }
}

// controlShip handles player input for ship movement
func (s *Spaceship) ShipControl() {
    // Horizontal movement
    if rl.IsKeyDown(rl.KeyLeft) && s.positon.X > MIN_X {
        s.positon.X -= SHIP_SPEED
    }
    if rl.IsKeyDown(rl.KeyRight) && s.positon.X < MAX_X {
        s.positon.X += SHIP_SPEED
    }

    // Vertical movement
    if rl.IsKeyDown(rl.KeyUp) && s.positon.Y > MIN_Y {
        s.positon.Y -= SHIP_SPEED
    }
    if rl.IsKeyDown(rl.KeyDown) && s.positon.Y < MAX_Y {
        s.positon.Y += SHIP_SPEED
    }
}

func (s *Spaceship) CheckCollisionWithRocks() bool {
    if (!s.isAlive) {
        return false
    }
    
    shipCenter, shipRadius := (&Projectile{position: s.positon}).GetBoundingCircle(rl.Vector2{X: 40, Y: 40})
    
    for _, rock := range Rocks {
        if rock.projectile.isActive {
            rockCenter, rockRadius := rock.projectile.GetBoundingCircle(rl.Vector2{X: 60, Y: 60})
            if CheckCircleCollision(shipCenter, shipRadius, rockCenter, rockRadius) {
                return true
            }
        }
    }
    return false
}

// Add method to handle player death
func (s *Spaceship) Die(gameState *GameState) {
    s.lives--
    s.isAlive = false
    gameState.isGameOver = true  // Immediately game over since only 1 life
}

// Remove respawn logic from Update since we don't respawn anymore
func (s *Spaceship) Update(deltaTime float32, gameState *GameState) {
    if s.isAlive {
        s.ShipControl()
    }
}


