package kamera

import "math"

// Vec2 for camera
type vec2 struct {
	X, Y float64
}

// add two vector
func (v vec2) Add(other vec2) vec2 {
	return vec2{v.X + other.X, v.Y + other.Y}
}

// sub returns this - other
func (v vec2) Sub(other vec2) vec2 {
	return vec2{v.X - other.X, v.Y - other.Y}
}

// lerp linearly interpolates between this and other vector.
func (v vec2) Lerp(other vec2, t float64) vec2 {
	return v.Scale(1.0 - t).Add(other.Scale(t))
}

// Scale scales vector
func (v vec2) Scale(factor float64) vec2 {
	return vec2{v.X * factor, v.Y * factor}
}

// Mult return this * other
func (v vec2) Mult(other vec2) vec2 {
	return vec2{v.X * other.X, v.Y * other.Y}
}

// Unit returns a normalized copy of this vector (unit vector).
func (v vec2) Unit() vec2 {
	// return v.Mult(1.0 / (v.Length() + math.SmallestNonzeroFloat64))
	return v.Scale(1.0 / (v.Mag() + 1e-50))
}

// Mag returns the magnitude of this vector
func (v vec2) Mag() float64 {
	return math.Sqrt(v.Dot(v))
}

// Dot returns dot product
func (v vec2) Dot(other vec2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// clamp returns f clamped to [low, high]
func clamp(f, low, high float64) float64 {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}
