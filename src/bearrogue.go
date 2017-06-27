package main

import (
	blt "bearlibterminal"
	"strconv"
)

const (
	WindowSizeX = 100
	WindowSizeY = 35
	Title = "BearRogue"
	Font = "fonts/UbuntuMono.ttf"
	FontSize = 24
)

var (
	// Global variables for now. This will be changed once we extrapolate this out to an Object
	playerX = 0
	playerY = 0
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

	blt.Color(blt.ColorFromName("white"))
	drawPlayer(playerX, playerY, "@")

	for {
		blt.Refresh()

		key := blt.Read()

		if key != blt.TK_CLOSE {
			handleInput(key)
			drawPlayer(playerX, playerY, "@")
		} else {
			break
		}

	}

	blt.Close()
}

// Handle basic character movement in the four main directions
func handleInput(key int) {
	switch key {
	case blt.TK_RIGHT:
		playerX ++
	case blt.TK_LEFT:
		playerX --
	case blt.TK_UP:
		playerY --
	case blt.TK_DOWN:
		playerY ++
	}
}

// Draw the player to the screen, at the given coordinates
func drawPlayer(x int, y int, symbol string) {
	blt.Layer(0)
	blt.ClearArea(0, 0, WindowSizeX, WindowSizeY)
	blt.Print(playerX, playerY, symbol)
}

// Centers text on the screen on the X and Y axis (taking string length into account)
func printCenteredText(text string) {
	x := WindowSizeX / 2 - (len(text) / 2)
	y := WindowSizeY / 2
	blt.Print(x, y, text)
}
