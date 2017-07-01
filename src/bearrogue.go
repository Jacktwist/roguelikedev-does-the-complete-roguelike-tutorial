package main

import (
	blt "bearlibterminal"
	"strconv"
	"entity"
	"gamemap"
)

const (
	WindowSizeX = 100
	WindowSizeY = 35
	MapWidth = WindowSizeX
	MapHeight = WindowSizeY
	Title = "BearRogue"
	Font = "fonts/UbuntuMono.ttf"
	FontSize = 24
)

var (
	player *entity.GameEntity
	entities []*entity.GameEntity
	gameMap *gamemap.Map
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
}
	
func main() {
	// Main game loop

	renderAll()

	for {
		blt.Refresh()

		key := blt.Read()

		// Clear each Entity off the screen
		for _, e := range entities {
			e.Clear()
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
		e.Draw()
	}
}

func renderMap() {
	// Render the game map. If a tile is blocked and blocks sight, draw a '#', if it is not blocked, and does not block
	// sight, draw a '.'
	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			if gameMap.Tiles[x][y].Blocked == true {
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
	renderMap()
	renderEntities()
}

