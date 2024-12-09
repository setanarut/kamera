package kamera

func lerp(start, end, t float64) float64 {
	return start + t*(end-start)
}
