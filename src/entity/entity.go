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
	Blocks bool
}

func (e *GameEntity) Move(dx int, dy int) {
	// Move the entity by the amount (dx, dy)
	e.X += dx
	e.Y += dy
}

func (e *GameEntity) Draw(mapX int, mapY int) {
	// Draw the entity to the screen
	blt.Layer(e.Layer)
	blt.Color(blt.ColorFromName(e.Color))
	blt.Print(mapX, mapY, e.Char)
}

func (e *GameEntity) Clear(mapX int, mapY int) {
	// Remove the entity from the screen
	blt.Layer(e.Layer)
	blt.Print(mapX, mapY, " ")
}
