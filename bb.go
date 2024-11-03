package kamera

// BB is Kamera's axis-aligned 2D bounding box type. (left, bottom, right, top)
type BB struct {
	L, B, R, T float64
}

// NewBBForExtents constructs a BB centered on a point with the given extents (half sizes).
func NewBBForExtents(centerX, centerY, halfWidth, halfHeight float64) BB {
	return BB{
		L: centerX - halfWidth,
		B: centerY - halfHeight,
		R: centerX + halfWidth,
		T: centerY + halfHeight,
	}
}

// ContainsPoint returns true if bb contains point.
func (bb BB) ContainsPoint(pointX, pointY float64) bool {
	return bb.L <= pointX && bb.R >= pointX && bb.B <= pointY && bb.T >= pointY
}

// Contains returns true if other lies completely within bb.
func (bb BB) Contains(other BB) bool {
	return bb.L <= other.L && bb.R >= other.R && bb.B <= other.B && bb.T >= other.T
}

// Intersects returns true if bb and other intersect.
func (bb BB) Intersects(other BB) bool {
	return bb.L <= other.R && other.L <= bb.R && bb.B <= other.T && other.B <= bb.T
}
