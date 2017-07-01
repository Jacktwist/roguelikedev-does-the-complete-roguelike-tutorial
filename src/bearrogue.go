package main

import (
	blt "bearlibterminal"
	"strconv"
	"entity"
)

const (
	WindowSizeX = 100
	WindowSizeY = 35
	Title = "BearRogue"
	Font = "fonts/UbuntuMono.ttf"
	FontSize = 24
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
}
	
func main() {
	// Main game loop

	player := entity.GameEntity{X: 1, Y: 1, Char: "@", Color: "red"}
	drawEntity(player)

	for {
		blt.Refresh()

		key := blt.Read()

		if key != blt.TK_CLOSE {
			handleInput(key, &player)
			drawEntity(player)
		} else {
			break
		}
	}

	blt.Close()
}

// Handle basic character movement in the four main directions
func handleInput(key int, player *entity.GameEntity) {
	switch key {
	case blt.TK_RIGHT:
		player.Move(1, 0)
	case blt.TK_LEFT:
		player.Move(-1, 0)
	case blt.TK_UP:
		player.Move(0, -1)
	case blt.TK_DOWN:
		player.Move(0, 1)
	}

	// Make sure the player cannot go outside the bounds of the window, for now
	// This will change when we later add camera controls
	if player.X> WindowSizeX - 1 {
		player.X = WindowSizeX - 1
	} else if player.X < 0 {
		player.X = 0
	}

	if player.Y > WindowSizeY - 1 {
		player.Y = WindowSizeY - 1
	} else if player.Y < 0 {
		player.Y = 0
	}
}

// Draw the player to the screen, at the given coordinates
func drawEntity(entity entity.GameEntity) {
	blt.Layer(0)
	blt.ClearArea(0, 0, WindowSizeX, WindowSizeY)
	blt.Color(blt.ColorFromName(entity.Color))
	blt.Print(entity.X, entity.Y, entity.Char)
}
