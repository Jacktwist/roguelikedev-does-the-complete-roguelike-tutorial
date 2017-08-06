package examinecursor

import (
	blt "bearlibterminal"
	"bearrogue/camera"
)

type XCursor struct {
	X int
	Y int
	Character string
	Layer int
}

func (c *XCursor) Move(dx, dy, maxX, maxY int, gameCamera *camera.GameCamera) {
	c.X += dx
	c.Y += dy

	if c.X < 0 {
		c.X = 0
	} else if c.X >= maxX {
		c.X = maxX - 2
	}

	if c.Y < 0 {
		c.Y = 0
	} else if c.Y >= maxY {
		c.Y = maxY - 2
	}
}

func (c *XCursor) Draw(gameCamera *camera.GameCamera) {
	blt.Layer(c.Layer)
	blt.Color(blt.ColorFromName("white"))
	cameraX, cameraY := gameCamera.ToCameraCoordinates(c.X, c.Y)
	blt.Print(cameraX, cameraY, c.Character)
}

func (c *XCursor) Clear(gameCamera *camera.GameCamera) {
	layer := blt.TK_LAYER
	blt.Layer(c.Layer)
	cameraX, cameraY := gameCamera.ToCameraCoordinates(c.X, c.Y)
	blt.Print(cameraX, cameraY, " ")
	blt.Layer(layer)
}