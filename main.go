package main

import (
	boids2 "Boid/boids"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"log"
)

func update(screen *ebiten.Image) error {
	if !ebiten.IsDrawingSkipped() {
		for _, boid := range boids2.Boids {
			screen.Set(int(boid.Position.X+1), int(boid.Position.Y), boid.Color)
			screen.Set(int(boid.Position.X-1), int(boid.Position.Y), boid.Color)
			screen.Set(int(boid.Position.X), int(boid.Position.Y-1), boid.Color)
			screen.Set(int(boid.Position.X), int(boid.Position.Y+1), boid.Color)
		}
	}
	return nil
}

func main() {
	for i := 0; i < boids2.BoidCount; i++ {
		println(fmt.Sprintf("\t\t %d  ", i))
		boids2.CreateBoid(i)
	}
	if err := ebiten.Run(update, boids2.ScreenWidth, boids2.ScreenHeight, 2, "Boids in a box"); err != nil {
		log.Fatal(err)
	}
}
