package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Projectile represents a moving object in the game (bullet/missile)
type Projectile struct {
	isActive bool         // Whether the projectile is currently in use
	position rl.Vector2   // Current position of the projectile
	speed    float32      // Movement speed of the projectile
	texture  rl.Texture2D // Visual representation of the projectile
}

func (p *Projectile) GetBoundingCircle(size rl.Vector2) (rl.Vector2, float32) {
	//calculate the center of projectile
	center := rl.NewVector2(
		p.position.X+size.X/2,
		p.position.Y+size.Y/2,
	)
	//calculate radius of the bounding circle
	radius := float32(math.Sqrt(float64((size.X/2)*(size.X/2) + (size.Y/2)*(size.Y/2))))

	return center, radius

}


func CheckCircleCollision(center1 rl.Vector2, radius1 float32, center2 rl.Vector2, radius2 float32) bool {
    // Calculate the distance between the centers of the two circles
    distance := rl.Vector2Distance(center1, center2)

    // Check if the distance is less than the sum of the radii
    return distance < (radius1 + radius2)
}

func CheckCircleRectangleCollision(center rl.Vector2, radius float32, rectPos rl.Vector2, rectSize rl.Vector2) bool {
	// Find closest point on rectangle to circle
	closestX := rl.Clamp(center.X, rectPos.X, rectPos.X+rectSize.X)
	closestY := rl.Clamp(center.Y, rectPos.Y, rectPos.Y+rectSize.Y)
	
	// Calculate distance between circle's center and closest point
	distanceX := center.X - closestX
	distanceY := center.Y - closestY
	
	// Check if distance is less than circle's radius
	distanceSquared := (distanceX * distanceX) + (distanceY * distanceY)
	return distanceSquared <= (radius * radius)
}