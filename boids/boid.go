package boids

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"
)

const (
	ScreenWidth, ScreenHeight = 640, 360
	BoidCount                 = 500
	viewRadius                = 13
	adjRate                   = 0.015
)

var (
	green       = color.RGBA{R: 10, G: 255, B: 50, A: 255}
	red         = color.RGBA{R: 240, G: 52, B: 52, A: 255}
	yellow      = color.RGBA{R: 255, G: 165, B: 0, A: 255}
	blue        = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	colors      = []color.RGBA{green, red, yellow, blue}
	Boids       [BoidCount]*Boid
	actualColor = -1
	boidMap     [ScreenWidth + 1][ScreenHeight + 1]int
)

func init() {
	for i, row := range boidMap {
		for j := range row {
			boidMap[i][j] = -1
		}
	}
}

type Boid struct {
	Position Vector2D
	Velocity Vector2D
	Id       int
	Color    color.RGBA
}

func (b *Boid) calcAcceleration() Vector2D {
	upper, lower := b.Position.AddV(viewRadius), b.Position.AddV(-viewRadius)
	avgVelocity := Vector2D{0, 0}
	count := 0.0
	for i := math.Max(lower.X, 0); i <= math.Min(upper.X, ScreenWidth); i++ {
		for j := math.Max(lower.Y, 0); j <= math.Min(upper.Y, ScreenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.Id {
				if dist := Boids[otherBoidId].Position.Distance(b.Position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(Boids[otherBoidId].Velocity)
				}
			}
		}
	}

	accel := Vector2D{0,0}
	if count > 0 {
		avgVelocity =  avgVelocity.DivisionV(count)
		accel = avgVelocity.Subtract(b.Velocity).MultiplyV(adjRate)
	}
	return accel
}

func (b *Boid) MoveOne() {
	b.Velocity = b.Velocity.Add(b.calcAcceleration()).Limit(-1, 1)
	boidMap[int(b.Position.X)][int(b.Position.Y)] = -1
	b.Position = b.Position.Add(b.Velocity)
	boidMap[int(b.Position.X)][int(b.Position.Y)] = b.Id
	next := b.Position.Add(b.Velocity)
	if next.X >= ScreenWidth || next.X < 0 {
		b.Velocity = Vector2D{-b.Velocity.X, b.Velocity.Y}
	}
	if next.Y >= ScreenHeight || next.Y < 0 {
		b.Velocity = Vector2D{b.Velocity.X, -b.Velocity.Y}
	}
}

func (b *Boid) Start() {
	for {
		b.MoveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func getColor() color.RGBA {
	if actualColor == len(colors)-1 {
		actualColor = -1
	}
	actualColor += 1
	return colors[actualColor]
}

func CreateBoid(bId int) {
	px := rand.Float64() * ScreenWidth
	py := rand.Float64() * ScreenHeight
	vx := (rand.Float64() * 2) - 1.0
	vy := (rand.Float64() * 2) - 1.0
	println(fmt.Sprintf("ID * %d * PX %f PY %f VX %f VY %f", bId, px, py, vx, vy))
	b := Boid{
		Position: Vector2D{X: px, Y: py},
		Velocity: Vector2D{X: vx, Y: vy},
		Id:       bId,
		Color:    getColor(),
	}
	boidMap[int(px)][int(py)] = bId
	Boids[bId] = &b
	go b.Start()
}
