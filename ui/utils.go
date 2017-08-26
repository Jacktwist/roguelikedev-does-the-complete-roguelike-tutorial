package ui

import (
	blt "bearlibterminal"
)

func ClearScreen(windowWidth, windowHeight int) {
	// Clear the entire screen, useful for showing various Menus etc
	blt.ClearArea(0, 0, windowWidth, windowHeight)

	for x := 0; x < windowWidth; x++ {
		for y := 0; y < windowHeight; y++ {
			// Clear both our primary layers, so we don't get any strange artifacts from one layer or the other getting
			// cleared.
			for i := 0; i <= 3; i++ {
				blt.Layer(i)
				blt.Print(x, y, " ")
			}
		}
	}
}

func MapBltKeyCodesToRunes(bltKeyCode int) rune {
	switch bltKeyCode {
	case blt.TK_A:
		return 'a'
	case blt.TK_B:
		return 'b'
	case blt.TK_C:
		return 'c'
	case blt.TK_D:
		return 'd'
	case blt.TK_E:
		return 'e'
	case blt.TK_F:
		return 'f'
	case blt.TK_G:
		return 'g'
	case blt.TK_H:
		return 'h'
	case blt.TK_I:
		return 'i'
	case blt.TK_J:
		return 'j'
	case blt.TK_K:
		return 'k'
	case blt.TK_L:
		return 'l'
	case blt.TK_M:
		return 'm'
	case blt.TK_N:
		return 'n'
	case blt.TK_O:
		return 'o'
	case blt.TK_P:
		return 'p'
	case blt.TK_Q:
		return 'q'
	case blt.TK_R:
		return 'r'
	case blt.TK_S:
		return 's'
	case blt.TK_T:
		return 't'
	case blt.TK_U:
		return 'u'
	case blt.TK_V:
		return 'v'
	case blt.TK_W:
		return 'w'
	case blt.TK_X:
		return 'x'
	case blt.TK_Y:
		return 'y'
	case blt.TK_Z:
		return 'z'
	}

	return ' '
}
