package main

type Vector struct {
	x int
	y int
}

func (vec Vector) neighbors() []Vector {
	return []Vector{
		vec.up(),
		vec.left(),
		vec.down(),
		vec.right(),
	}
}

func (vec Vector) down() Vector {
	vec.y += 1
	return vec
}
func (vec Vector) up() Vector {
	vec.y -= 1
	return vec
}
func (vec Vector) right() Vector {
	vec.x += 1
	return vec
}
func (vec Vector) left() Vector {
	vec.x -= 1
	return vec
}
func (vec Vector) inBounds(xBounds int, yBounds int) bool {
	if vec.x < 0 || vec.y < 0 || vec.x >= xBounds || vec.y >= yBounds {
		return false
	}
	return true
}
