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
	blt.Color(blt.ColorFromName("darker green"))
	printCenteredText("Hello, World!")
	blt.Refresh()

	for blt.Read() != blt.TK_CLOSE {
		// Do nothing for now
	}

	blt.Close()
}

// Centers text on the screen on the X and Y axis (taking string length into account)
func printCenteredText(text string) {
	x := WindowSizeX / 2 - (len(text) / 2)
	y := WindowSizeY / 2
	blt.Print(x, y, text)
}
