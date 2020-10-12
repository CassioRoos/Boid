package boids

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	ScreenWidth, ScreenHeight = 640, 360
	BoidCount                 = 500
	viewRadius                = 10
	adjRate                   = 0.015
)

var (
	green       = color.RGBA{R: 10, G: 255, B: 50, A: 255}
	red         = color.RGBA{R: 240, G: 52, B: 52, A: 255}
	yellow      = color.RGBA{R: 255, G: 165, B: 0, A: 255}
	blue        = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	white       = color.RGBA{R: 255, G: 255, B: 255, A: 1}
	fuchsia     = color.RGBA{R: 255, G: 0, B: 255, A: 1}
	colors      = []color.RGBA{white, green, red, yellow, blue, fuchsia}
	Boids       [BoidCount]*Boid
	actualColor = -1
	boidMap     [ScreenWidth + 1][ScreenHeight + 1]int
	rwLock      = sync.RWMutex{}
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
	avgPosition, avgVelocity, separation := Vector2D{0, 0}, Vector2D{0, 0}, Vector2D{0, 0}
	count := 0.0
	rwLock.RLock()
	for i := math.Max(lower.X, 0); i <= math.Min(upper.X, ScreenWidth); i++ {
		for j := math.Max(lower.Y, 0); j <= math.Min(upper.Y, ScreenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.Id {
				if dist := Boids[otherBoidId].Position.Distance(b.Position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(Boids[otherBoidId].Velocity)
					avgPosition = avgPosition.Add(Boids[otherBoidId].Position)
					separation = separation.Add(b.Position.Subtract(Boids[otherBoidId].Position).DivisionV(dist))

				}
			}
		}
	}
	rwLock.RUnlock()
	accel := Vector2D{b.borderBounce(b.Position.X, ScreenWidth), b.borderBounce(b.Position.Y, ScreenHeight)}
	if count > 0 {
		avgVelocity, avgPosition = avgVelocity.DivisionV(count), avgPosition.DivisionV(count)
		accelAligment := avgVelocity.Subtract(b.Velocity).MultiplyV(adjRate)
		accelCohesion := avgPosition.Subtract(b.Position).MultiplyV(adjRate)
		accelSeparation := separation.MultiplyV(adjRate * 2)
		accel = accel.Add(accelAligment).Add(accelCohesion).Add(accelSeparation)
	}
	return accel
}

func (b *Boid) MoveOne() {
	calcAcceleration := b.calcAcceleration()
	rwLock.Lock()
	b.Velocity = b.Velocity.Add(calcAcceleration).Limit(-1, 1)
	boidMap[int(b.Position.X)][int(b.Position.Y)] = -1
	b.Position = b.Position.Add(b.Velocity)
	boidMap[int(b.Position.X)][int(b.Position.Y)] = b.Id
	rwLock.Unlock()
}

func (b *Boid) Start() {
	for {
		b.MoveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func (b *Boid) borderBounce(pos, maxBorderPos float64) float64 {
	if pos < viewRadius {
		return 1 / pos
	} else if pos > maxBorderPos-viewRadius {
		return 1 / (pos - maxBorderPos)
	}
	return 0
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
