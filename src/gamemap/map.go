package gamemap

type Tile struct {
	Blocked bool
	Blocks_sight bool
}

type Map struct {
	Width  int
	Height int
	Tiles  [][]*Tile
}

func (m *Map) InitializeMap() {
	// Set up a map where all the border (edge) Tiles are walls (block movement, and sight)
	// This is just a test method, we will build maps more dynamically in the future.
	m.Tiles = make([][]*Tile, m.Width)
	for i := range m.Tiles {
		m.Tiles[i] = make([]*Tile, m.Height)
	}

	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if x == 0 || x == m.Width- 1 || y == 0 || y == m.Height- 1 {
				m.Tiles[x][y] = &Tile{true, true}
			} else {
				m.Tiles[x][y] = &Tile{false, false}
			}
		}
	}
}

func (m *Map) IsBlocked(x int, y int) bool {
	// Check to see if the provided coordinates contain a blocked tile
	if m.Tiles[x][y].Blocked {
		return true
	} else {
		return false
	}
}
