package model

// Vector is a simple vector struct
type Vector struct {
	X float64
	Y float64
	Z float64
}

// NewVector is a factory function creating instance of Vector
func NewVector() Vector {
	return Vector{
		X: 0,
		Y: 0,
		Z: 0,
	}
}
