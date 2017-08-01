package gamemap

func (m *Map) GenerateArena() {
	// Generates a large, empty room, with walls ringing the outside edges
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if x == 0 || x == m.Width-1 || y == 0 || y == m.Height-1 {
				m.Tiles[x][y] = &Tile{true, true, false, false, false, x, y}
			} else {
				m.Tiles[x][y] = &Tile{false, false, false, false, false, x, y}
			}
		}
	}
}
