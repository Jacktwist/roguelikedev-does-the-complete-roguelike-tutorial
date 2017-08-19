package examinecursor

import (
	blt "bearlibterminal"
	"bearrogue/camera"
)

type XCursor struct {
	X         int
	Y         int
	Character string
	Layer     int
}

func (c *XCursor) Move(dx, dy, maxX, maxY int, gameCamera *camera.GameCamera) {

	oldCX, oldCY := c.X, c.Y

	c.X += dx
	c.Y += dy

	cameraX, cameraY := gameCamera.ToCameraCoordinates(c.X, c.Y)

	// Check if the cursor is outside the camera bounds. If it is, set it back to it original position.
	if cameraX < 0 {
		c.X = oldCX
	} else if cameraX >= maxX {
		c.X = oldCX
	}

	if cameraY < 0 {
		c.Y = oldCY
	} else if cameraY >= maxY {
		c.Y = oldCY
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
