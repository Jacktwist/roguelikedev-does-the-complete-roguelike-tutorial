package gamemap

import (
	"math/rand"
	"sort"
	"time"
)

type BySize [][]*Tile

func (s BySize) Len() int {
	return len(s)
}

func (s BySize) Swap(i, j int) {
	s[i], s[j]= s[j], s[i]
}

func (s BySize) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

type Tile struct {
	Blocked bool
	Blocks_sight bool
	Visited bool
	Explored bool
	Visible bool
	X int
	Y int
}

func (t *Tile) IsWall() bool {
	if t.Blocks_sight && t.Blocked {
		return true
	} else {
		return false
	}
}

type Map struct {
	Width  int
	Height int
	Tiles  [][]*Tile
}

func (m *Map) InitializeMap() {
	// Initialize a two dimensional array that will represent the current game map (of dimensions Width x Height)
	m.Tiles = make([][]*Tile, m.Width)
	for i := range m.Tiles {
		m.Tiles[i] = make([]*Tile, m.Height)
	}

	// Set a seed for procedural generation
	rand.Seed( time.Now().UTC().UnixNano())
}

func (m *Map) IsBlocked(x int, y int) bool {
	// Check to see if the provided coordinates contain a blocked tile
	if m.Tiles[x][y].Blocked {
		return true
	} else {
		return false
	}
}

func (m *Map) GenerateArena() {
	// Generates a large, empty room, with walls ringing the outside edges
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if x == 0 || x == m.Width- 1 || y == 0 || y == m.Height- 1 {
				m.Tiles[x][y] = &Tile{true, true, false, false, false, x, y}
			} else {
				m.Tiles[x][y] = &Tile{false, false, false, false, false, x, y}
			}
		}
	}
}

func (m *Map) GenerateCavern() (int, int) {

	// Step 1: Fill the map space with a random assortment of walls and floors. This uses a roughly 40/60 ratio in favor
	// of floors, as I've found that to produce the nicest results.
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			state := rand.Intn(100)
			if state < 50 {
				m.Tiles[x][y] = &Tile{true, true, false, false, false, x, y}
			} else {
				m.Tiles[x][y] = &Tile{false, false, false, false, false, x, y}
			}
		}
	}

	// Step 2: Decide what should remain as walls. If four or more of a tiles immediate (within 1 space) neighbors are
	// walls, then make that tile a wall. If 2 or less of the tiles next closest (2 spaces away) neighbors are walls,
	// then make that tile a wall. Any other scenario, and the tile will become (or stay) a floor tile.
	// Make several passes on this to help smooth out the walls of the cave.
	for i := 0; i < 5; i++ {
		for x := 0; x < m.Width; x++ {
			for y := 0; y < m.Height - 1; y++ {
				wallOneAway := m.countWallsNStepsAway(1, x, y)

				wallTwoAway := m.countWallsNStepsAway(2, x, y)

				if wallOneAway >= 5 || wallTwoAway <= 2 {
					m.Tiles[x][y].Blocked = true
					m.Tiles[x][y].Blocks_sight = true
				} else {
					m.Tiles[x][y].Blocked = false
					m.Tiles[x][y].Blocks_sight = false
				}
			}
		}
	}

	// Step 3: Make a few more passes, smoothing further, and removing any small or single tile, unattached walls.
	for i := 0; i < 5; i++ {
		for x := 0; x < m.Width; x++ {
			for y := 0; y < m.Height - 1; y++ {
				wallOneAway := m.countWallsNStepsAway(1, x, y)

				if wallOneAway >= 5 {
					m.Tiles[x][y].Blocked = true
					m.Tiles[x][y].Blocks_sight = true
				} else {
					m.Tiles[x][y].Blocked = false
					m.Tiles[x][y].Blocks_sight = false
				}
			}
		}
	}

	// Step 4: Seal up the edges of the map, so the player, and the following flood fill passes, cannot go beyond the
	// intended game area
	for x := 0; x < m.Width ; x++ {
		for y := 0; y < m.Height; y++ {
			if x == 0 || x == m.Width - 1 || y == 0 || y == m.Height - 1 {
				m.Tiles[x][y].Blocked = true
				m.Tiles[x][y].Blocks_sight = true
			}
		}
	}

	// Step 5: Flood fill. This will find each individual cavern in the cave system, and add them to a list. It will
	// then find the largest one, and will make that as the main play area. The smaller caverns will be filled in.
	// In the future, it might make sense to tunnel between caverns, and apply a few more smoothing passes, to make
	// larger, more realistic caverns.

	var cavern []*Tile
	var totalCavernArea []*Tile
	var caverns [][]*Tile
	var tile *Tile
	var node *Tile

	for x := 0; x < m.Width - 1; x++ {
		for y := 0; y < m.Height - 1; y++ {
			tile = m.Tiles[x][y]

			// If the current tile is a wall, or has already been visited, ignore it and move on
			if !tile.Visited && !tile.IsWall() {
				// This is a non-wall, unvisited tile
				cavern = append(cavern, m.Tiles[x][y])

				for len(cavern) > 0 {
					// While the current node tile has valid neighbors, keep looking for more valid neighbors off of
					// each one
					node = cavern[len(cavern)-1]
					cavern = cavern[:len(cavern)-1]

					if !node.Visited && !node.IsWall() {
						// Mark the node as visited, and add it to the cavern area for this cavern
						node.Visited = true
						totalCavernArea = append(totalCavernArea, node)

						// Add the tile to the west, if valid
						if node.X - 1 > 0 && !m.Tiles[node.X -1][node.Y].IsWall() {
							cavern = append(cavern, m.Tiles[node.X -1][node.Y])
						}

						// Add the tile to east, if valid
						if node.X + 1 < m.Width && !m.Tiles[node.X + 1][node.Y].IsWall() {
							cavern = append(cavern, m.Tiles[node.X + 1][node.Y])
						}

						// Add the tile to north, if valid
						if node.Y - 1 > 0 && !m.Tiles[node.X][node.Y - 1].IsWall() {
							cavern = append(cavern, m.Tiles[node.X][node.Y - 1])
						}

						// Add the tile to south, if valid
						if node.Y + 1 < m.Height && !m.Tiles[node.X][node.Y + 1].IsWall() {
							cavern = append(cavern, m.Tiles[node.X][node.Y + 1])
						}
					}
				}

				// All non-wall tiles have been found for the current cavern, add it to the list, and start looking for
				// the next one
				caverns = append(caverns, totalCavernArea)
				totalCavernArea = nil
			} else {
				tile.Visited = true
			}
		}
	}

	// Sort the caverns slice by size. This will make the largest cavern last, which will then be removed from the list.
	// Then, fill in any remaining caverns (aside from the main one). This will ensure that there are no areas on the
	// map that the player cannot reach.
	sort.Sort(BySize(caverns))
	mainCave := caverns[len(caverns) - 1]
	caverns = caverns[:len(caverns) - 1]

	for i := 0; i < len(caverns); i++ {
		for j := 0; j < len(caverns[i]); j++ {
			caverns[i][j].Blocked = true
			caverns[i][j].Blocks_sight = true
		}
	}

	// Finally, choose a starting position for the player within the newly created cave
	pos := rand.Int() % len(mainCave)
	return mainCave[pos].X, mainCave[pos].Y

	return 0, 0
}

func (m *Map) countWallsNStepsAway(n int, x int, y int) int {
	// Return the number of wall tiles that are within n spaces of the given tile
	wallCount := 0

	for r := -n; r <= n; r++ {
		for c := -n; c <= n; c++ {
			if x + r >= m.Width || x + r <= 0 || y + c >= m.Height || y + c <= 0 {
				// Check if the current coordinates would be off the map. Off map coordinates count as a wall.
				wallCount ++
			} else if m.Tiles[x + r][y + c].Blocked && m.Tiles[x + r][y + c].Blocks_sight {
				wallCount ++
			}
		}
	}

	return wallCount
}