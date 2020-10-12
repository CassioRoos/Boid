package boids

import "math"

type Vector2D struct {
	X float64
	Y float64
}

func (v1 Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func (v1 Vector2D) Subtract(v2 Vector2D) Vector2D {
	return Vector2D{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func (v1 Vector2D) Multiply(v2 Vector2D) Vector2D {
	return Vector2D{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func (v1 Vector2D) AddV(v float64) Vector2D {
	return Vector2D{X: v1.X + v, Y: v1.Y + v}
}

func (v1 Vector2D) MultiplyV(v float64) Vector2D {
	return Vector2D{X: v1.X * v, Y: v1.Y * v}
}

func (v1 Vector2D) DivisionV(v float64) Vector2D {
	return Vector2D{X: v1.X / v, Y: v1.Y / v}
}

func (v1 Vector2D) Limit(lower, upper float64) Vector2D {
	return Vector2D{math.Min(math.Max(v1.X, lower), upper),
		math.Min(math.Max(v1.Y, lower), upper)}
}

func (v1 Vector2D) Distance(v2 Vector2D) float64 {
	return math.Sqrt(math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2))
}
