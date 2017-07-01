package entity

import (
	blt "bearlibterminal"
)

type GameEntity struct {
	X int
	Y int
	Layer int
	Char string
	Color string
}

func (e *GameEntity) Move(dx int, dy int) {
	// Move the entity by the amount (dx, dy)
	e.X += dx
	e.Y += dy
}

func (e *GameEntity) Draw() {
	// Draw the entity to the screen
	blt.Layer(e.Layer)
	blt.Color(blt.ColorFromName(e.Color))
	blt.Print(e.X, e.Y, e.Char)
}

func (e *GameEntity) Clear() {
	// Remove the entity from the screen
	blt.Layer(e.Layer)
	blt.Print(e.X, e.Y, " ")
}
