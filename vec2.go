package kamera

import "math"

// vec2 for camera
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

// smoothDamp gradually changes a value towards a desired goal over time.
func smoothDamp(current, target vec2, currentVelocity *vec2, smoothTime, maxSpeed float64) vec2 {
	smoothTime = math.Max(0.0001, smoothTime)
	omega := 2.0 / smoothTime
	x := omega * 0.016666666666666666
	exp := 1.0 / (1.0 + x + 0.48*x*x + 0.235*x*x*x)
	change := current.Sub(target)
	originalTo := target
	maxChange := maxSpeed * smoothTime
	maxChangeSq := maxChange * maxChange
	sqDist := change.Dot(change)
	if sqDist > maxChangeSq {
		mag := math.Sqrt(sqDist)
		change = change.Scale(maxChange / mag)
	}
	target = current.Sub(change)
	temp := (currentVelocity.Add(vec2{change.X * omega, change.Y * omega})).Scale(0.016666666666666666)
	*currentVelocity = currentVelocity.Sub(vec2{temp.X * omega, temp.Y * omega}).Scale(exp)
	output := target.Add(change.Add(temp).Scale(exp))
	origMinusCurrent := originalTo.Sub(current)
	outMinusOrig := output.Sub(originalTo)
	if origMinusCurrent.Dot(outMinusOrig) > 0 {
		output = originalTo
		*currentVelocity = output.Sub(originalTo).Scale(1.0 / 0.016666666666666666)
	}
	return output
}
