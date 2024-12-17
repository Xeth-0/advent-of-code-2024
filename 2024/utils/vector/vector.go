package vector

import "math"

// Was writing this in each day's code that required it, but since it seems like
// basically every puzzle just has to be a 2D grid puzzle...
type Vector struct {
	X int
	Y int
}

func (vec *Vector) SetX(x int) {
	vec.X = x
}

func (vec *Vector) SetY(y int) {
	vec.Y = y
}

func (vec *Vector) Set(x int, y int) {
	vec.X = x
	vec.Y = y
}

func (vec *Vector) Equals(vec2 Vector) bool {
	return vec.X == vec2.X && vec.Y == vec2.Y
}

func (vec Vector) Neighbors() []Vector {
	return []Vector{
		vec.Up(),
		vec.Left(),
		vec.Down(),
		vec.Right(),
	}
}

func (vec Vector) Down() Vector {
	vec.Y += 1
	return vec
}
func (vec Vector) Up() Vector {
	vec.Y -= 1
	return vec
}
func (vec Vector) Right() Vector {
	vec.X += 1
	return vec
}
func (vec Vector) Left() Vector {
	vec.X -= 1
	return vec
}
func (vec Vector) InBounds(xBounds int, yBounds int) bool {
	if vec.X < 0 || vec.Y < 0 || vec.X >= xBounds || vec.Y >= yBounds {
		return false
	}
	return true
}

func (vec Vector) Add(vec2 Vector) Vector {
	vec.X += vec2.X
	vec.Y += vec2.Y
	return vec
}

func (vec Vector) Distance(vec2 Vector) int {
	// d2 := float64((vec.X - vec2.X) ^ 2 + (vec.Y - vec2.Y) ^ 2)
	// return int(math.Sqrt(d2))

	d2 := math.Abs(float64(vec2.X-vec.X)) + math.Abs(float64(vec2.Y-vec.Y))
	return int(d2)
}
