package main

import (
	blt "bearlibterminal"
	"camera"
	"entity"
	"gamemap"
	"strconv"
)

const (
	WindowSizeX = 100
	WindowSizeY = 35
	MapWidth = 200
	MapHeight = 200
	Title = "BearRogue"
	Font = "fonts/UbuntuMono.ttf"
	FontSize = 24
)

var (
	player *entity.GameEntity
	entities []*entity.GameEntity
	gameMap *gamemap.Map
	gameCamera *camera.GameCamera
)

func init() {
	blt.Open()

	// BearLibTerminal uses configuration strings to set itself up, so we need to build these strings here
	// First set up the string for window properties (size and title)
	size := "size=" + strconv.Itoa(WindowSizeX) + "x" + strconv.Itoa(WindowSizeY)
	title := "title='" + Title + "'"
	window := "window: " + size + "," + title

	// Next, setup the font config string
	fontSize := "size=" + strconv.Itoa(FontSize)
	font := "font: " + Font + ", " + fontSize

	// Now, put it all together
	blt.Set(window + "; " + font)
	blt.Clear()

	// Create a player Entity and an NPC entity, and add them to our slice of Entities
	player = &entity.GameEntity{X: 1, Y: 1, Layer: 1, Char: "@", Color: "white"}
	npc := &entity.GameEntity{X: 10, Y: 10, Layer: 0, Char: "N", Color: "red"}
	entities = append(entities, player, npc)

	// Create a GameMap, and initialize it
	gameMap = &gamemap.Map{Width: MapWidth, Height: MapHeight}
	gameMap.InitializeMap()

	// Initialize a camera object
	gameCamera = &camera.GameCamera{X: 1, Y:1, Width: WindowSizeX, Height: WindowSizeY}
}
	
func main() {
	// Main game loop

	renderAll()

	for {
		blt.Refresh()

		key := blt.Read()

		// Clear each Entity off the screen
		for _, e := range entities {
			mapX, mapY := gameCamera.ToCameraCoordinates(e.X, e.Y)
			e.Clear(mapX, mapY)
		}

		if key != blt.TK_CLOSE {
			handleInput(key, player)
		} else {
			break
		}
		renderAll()
	}

	blt.Close()
}

func handleInput(key int, player *entity.GameEntity) {
	// Handle basic character movement in the four main directions

	var (
		dx, dy int
	)

	switch key {
	case blt.TK_RIGHT:
		dx, dy = 1, 0
	case blt.TK_LEFT:
		dx, dy = -1, 0
	case blt.TK_UP:
		dx, dy = 0, -1
	case blt.TK_DOWN:
		dx, dy = 0, 1
	}

	// Check to ensure that the tile the player is trying to move in to is a valid move (not blocked)
	if !gameMap.IsBlocked(player.X + dx, player.Y + dy) {
		player.Move(dx, dy)
	}
}

func renderEntities() {
	// Draw every Entity present in the game. This gets called on each iteration of the game loop.
	for _, e := range entities {
		mapX, mapY := gameCamera.ToCameraCoordinates(e.X, e.Y)
		e.Draw(mapX, mapY)
	}
}

func renderMap() {
	// Render the game map. If a tile is blocked and blocks sight, draw a '#', if it is not blocked, and does not block
	// sight, draw a '.'
	for y := 0; y < gameCamera.Height; y++ {
		for x := 0; x < gameCamera.Width; x++ {
			mapX, mapY := gameCamera.X + x, gameCamera.Y + y
			if gameMap.Tiles[mapX][mapY].Blocked == true {
				blt.Color(blt.ColorFromName("gray"))
				blt.Print(x, y, "#")
			} else {
				blt.Color(blt.ColorFromName("brown"))
				blt.Print(x, y, ".")
			}
		}
	}
}

func renderAll() {
	// Convenience function to render all entities, followed by rendering the game map

	// Before anything is rendered, update the camera position, so it is centered (if possible) on the player
	// Only things wintin the cameras viewport will be drawn to the screen
	gameCamera.MoveCamera(player.X, player.Y, MapWidth, MapHeight)

	renderMap()
	renderEntities()
}

