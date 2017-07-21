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
	Name string
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

func GetBlockingEntitiesAtLocation(entities []*GameEntity, destinationX, destinationY int) *GameEntity {
	// Return any entities that are at the destination location which would block movement
	for _, e := range entities {
		if e.Blocks && e.X == destinationX && e.Y == destinationY {
			return e
		}
	}
	return nil
}
