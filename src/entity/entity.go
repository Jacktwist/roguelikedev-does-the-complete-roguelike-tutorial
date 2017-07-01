package entity

type GameEntity struct {
	X int
	Y int
	Char string
	Color string
}

func (e *GameEntity) Move(dx int, dy int) {
	// Move the entity by the amount (dx, dy)
	e.X += dx
	e.Y += dy
}
