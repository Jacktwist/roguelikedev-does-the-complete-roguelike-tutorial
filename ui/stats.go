package ui

import (
	blt "bearlibterminal"
	"strconv"
	"math"
)

func PrintBasicCharacterInfo(name string, viewAreaX int) {
	// Print basic information about the player character to the info sidebar

	startX := viewAreaX
	startY := 1

	blt.Print(startX, startY, name)
	blt.Print(startX, startY + 1, "Unremarkable Human")
	blt.Print(startX, startY + 2, "\n")

}

func PrintStats(hp, maxHp, viewAreaX int) {
	startX := viewAreaX
	startY := 3

	printHpBar(hp, maxHp, startX, startY)
	blt.Print(startX, startY + 1, "ST: (20/20) [color=yellow]==========[/color]")
	blt.Print(startX, startY + 2, "MG: (20/20) [color=blue]==========[/color]")
}

func printHpBar(hp, maxHp, startX, startY int) {

	numericRepresentation := ""
	if hp < 10 {
		// Add some padding for single digits, so things line up
		numericRepresentation = "( " + strconv.Itoa(hp) + "/" + strconv.Itoa(maxHp) + ")"
	} else {
		numericRepresentation = "(" + strconv.Itoa(hp) + "/" + strconv.Itoa(maxHp) + ")"
	}

	// Figure out how many health pips to display, based on the players current health
	percent := float64(hp) / float64(maxHp)
	pips := int(round(percent * 10))

	healthBar := ""
	for i := 0; i < 10; i++ {
		if i < pips {
			healthBar += "="
		} else {
			healthBar += "-"
		}
	}

	blt.Print(startX, startY, "HP: " + numericRepresentation +" [color=red]" + healthBar + "[/color]")
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}
