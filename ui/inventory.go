package ui

import (
	blt "bearlibterminal"
	"strconv"
)

func printHeader(title string, invMax, invUsed int) {
	blt.Print(1, 1, title+" ("+strconv.Itoa(invUsed)+"/"+strconv.Itoa(invMax)+")")
	blt.Print(1, 2, "--------------------")
}

func printInventoryItems(items map[string]int) {
	y := 3

	if len(items) > 0 {
		for k, v := range items {
			blt.Print(1, y, k+" x"+strconv.Itoa(v))
			y++
		}
	} else {
		blt.Print(1, y, "You are not carrying anything...")
	}
}

func DisplayInventory(title string, invMax, invUsed int, items map[string]int) {
	printHeader(title, invMax, invUsed)
	printInventoryItems(items)
}

func DisplayInformationScreen(title, shortDescription, longDescription string, occurences int, windowHeight int) {
	blt.Print(1, 1, title)
	blt.Print(1, 3, shortDescription)

	if longDescription != "" {
		blt.Print(1, 6, longDescription)
	}

	blt.Print(1, 10, "You have "+strconv.Itoa(occurences)+" of these.")

	blt.Print(1, windowHeight-1, "[color=light blue]Actions available: (d)rop this item[/color]")
}
