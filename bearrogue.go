package main

import (
	blt "bearlibterminal"
	"bearrogue/camera"
	"bearrogue/ecs"
	"bearrogue/fov"
	"bearrogue/gamemap"
	"strconv"
	"bearrogue/ui"
	"math/rand"
	"fmt"
)

const (
	WindowSizeX = 100
	WindowSizeY = 35
	ViewAreaX = 100
	ViewAreaY = 30
	MapWidth = 100
	MapHeight = 35
	Title = "BearRogue"
	Font = "fonts/UbuntuMono.ttf"
	FontSize = 24
	PlayerTurn = iota
	MobTurn = iota
)

var (
	player *ecs.GameEntity
	entities []*ecs.GameEntity
	gameMap *gamemap.Map
	gameCamera *camera.GameCamera
	fieldOfView *fov.FieldOfVision
	gameTurn int
	messageLog ui.MessageLog
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

	// Create a player Entity, and add them to our slice of Entities
	player = &ecs.GameEntity{}
	player.SetupGameEntity()
	player.AddComponent("player", ecs.PlayerComponent{})
	player.AddComponent("position", ecs.PositionComponent{X: 0, Y: 0})
	player.AddComponent("appearance", ecs.AppearanceComponent{Color: "white", Character: "@", Layer: 1, Name: "Player"})
	player.AddComponent("movement", ecs.MovementComponent{})
	player.AddComponent("controllable", ecs.ControllableComponent{})
	player.AddComponent("attacker", ecs.AttackerComponent{Attack: 5, Defense: 5})
	player.AddComponent("hitpoints", ecs.HitPointComponent{Hp: 20, MaxHP: 20})
	player.AddComponent("block", ecs.BlockingComponent{})
	player.AddComponent("killable", ecs.KillableComponent{Name: "Here lies", Character: "%", Color: "dark red"})

	entities = append(entities, player)

	// Create a GameMap, and initialize it (and set the player position within it, for now)
	gameMap = &gamemap.Map{Width: MapWidth, Height: MapHeight}
	gameMap.InitializeMap()

	playerX, playerY, mapEntities := GenerateAndPopulateCavern()

	if player.HasComponent("position") {
		positionComponent, _ := player.Components["position"].(ecs.PositionComponent)
		positionComponent.X = playerX
		positionComponent.Y = playerY
		player.RemoveComponent("position")
		player.AddComponent("position", positionComponent)
		player.Print()
	}

	entities = append(entities, mapEntities...)

	// Set the current turn to the player, so they may act first
	gameTurn = PlayerTurn

	// Initialize a camera object
	gameCamera = &camera.GameCamera{X: 1, Y:1, Width: ViewAreaX, Height: ViewAreaY}

	// Initialize a FoV object
	fieldOfView = &fov.FieldOfVision{}
	fieldOfView.Initialize()
	fieldOfView.SetTorchRadius(6)

	// Set up the messageLog, and output a "welcome" message
	messageLog = ui.MessageLog{MaxLength: 100}
	messageLog.InitMessages()
}
	
func main() {
	// Main game loop

	messageLog.SendMessage("You find yourself in the caverns of eternal sadness...you start to feel a little more sad.")
	messageLog.PrintMessages(ViewAreaY, WindowSizeX, WindowSizeY)
	renderMap()
	ecs.SystemRender(entities, gameCamera, gameMap)

	for {
		blt.Refresh()

		key := blt.Read()

		// Clear each Entity off the screen
		ecs.SystemClear(entities, gameCamera)

		if key != blt.TK_CLOSE {
			if gameTurn == PlayerTurn {
				if player.HasComponents([]string{"movement", "controllable", "position"}) {
					handleInput(key, player)
				}
			}
		} else {
			break
		}

		if gameTurn == MobTurn {
			var newEntities []*ecs.GameEntity
			for _, e := range entities {
				if e != nil {
					if !e.HasComponent("player") {
						ecs.SystemMovement(e, 0, 0, entities, gameMap, &messageLog)
						newEntities = append(newEntities, ecs.SystemReproduce(e, entities, gameMap, &messageLog))
					}
				}
			}
			entities = append(entities, newEntities...)
			gameTurn = PlayerTurn
		}

		renderMap()
		ecs.SystemRender(entities, gameCamera, gameMap)
		messageLog.PrintMessages(ViewAreaY, WindowSizeX, WindowSizeY)
	}

	blt.Close()
}

func handleInput(key int, entity *ecs.GameEntity) {
	// Handle basic character movement in the four main directions, plus diagonals (and vim keys)

	var (
		dx, dy int
	)

	switch key {
	case blt.TK_RIGHT, blt.TK_L:
		dx, dy = 1, 0
	case blt.TK_LEFT, blt.TK_H:
		dx, dy = -1, 0
	case blt.TK_UP, blt.TK_K:
		dx, dy = 0, -1
	case blt.TK_DOWN, blt.TK_J:
		dx, dy = 0, 1
	case blt.TK_Y:
		dx, dy = -1, -1
	case blt.TK_U:
		dx, dy = 1, -1
	case blt.TK_B:
		dx, dy = -1, 1
	case blt.TK_N:
		dx, dy = 1, 1
	}

	// Fire off the movement system
	ecs.SystemMovement(entity, dx, dy, entities, gameMap, &messageLog)

	// Switch the game turn to the Mobs turn
	gameTurn = MobTurn
}

func renderMap() {
	// Render the game map. If a tile is blocked and blocks sight, draw a '#', if it is not blocked, and does not block
	// sight, draw a '.'

	// First, set the entire map to not visible. We'll decide what is visible based on the torch radius.
	// In the process, clear every Tile on the map as well
	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			gameMap.Tiles[x][y].Visible = false
			blt.Print(x, y, " ")
		}
	}

	positionComponent, posOk := player.Components["position"].(ecs.PositionComponent)

	if posOk {
		gameCamera.MoveCamera(positionComponent.X, positionComponent.Y, MapWidth, MapHeight)

		// Next figure out what is visible to the player, and what is not.
		fieldOfView.RayCast(positionComponent.X, positionComponent.Y, gameMap)
	}

	// Now draw each tile that should appear on the screen, if its visible, or explored
	for x := 0; x < gameCamera.Width; x++ {
		for y := 0; y < gameCamera.Height; y++ {
			mapX, mapY := gameCamera.X + x, gameCamera.Y + y

			if gameMap.Tiles[mapX][mapY].Visible {
				if gameMap.Tiles[mapX][mapY].IsWall() {
					blt.Color(blt.ColorFromName("white"))
					blt.Print(x, y, "#")
				} else {
					blt.Color(blt.ColorFromName("white"))
					blt.Print(x, y, ".")
				}
			} else if gameMap.Tiles[mapX][mapY].Explored {
				if gameMap.Tiles[mapX][mapY].IsWall() {
					blt.Color(blt.ColorFromName("gray"))
					blt.Print(x, y, "#")
				} else {
					blt.Color(blt.ColorFromName("gray"))
					blt.Print(x, y, ".")
				}
			}
		}
	}
}

/* Generator functions */
func GenerateAndPopulateCavern() (int, int, []*ecs.GameEntity) {
	gameMap := gameMap.GenerateCavern()

	pos := rand.Int() % len(gameMap)
	playerX, playerY := gameMap[pos].X, gameMap[pos].Y

	entities := populateCavern(gameMap)

	return playerX, playerY, entities
}

func populateCavern(mainCave []*gamemap.Tile) []*ecs.GameEntity {
	// Randomly sprinkle some Orcs, Trolls, and Goblins around the newly created cavern
	var entities []*ecs.GameEntity
	var createdEntity *ecs.GameEntity

	for i := 0; i < 10; i++ {
		x := 0
		y := 0
		locationFound := false
		for j := 0; j <= 50; j++ {
			// Attempt to find a clear location to create a mob (ecs for now)
			pos := rand.Int() % len(mainCave)
			x = mainCave[pos].X
			y = mainCave[pos].Y
			if ecs.GetBlockingEntitiesAtLocation(entities, x, y) == nil {
				locationFound = true
				break
			}
		}

		if locationFound {
			chance := rand.Intn(100)
			if chance <= 1 {
				// Create a Troll
				createdEntity = &ecs.GameEntity{}
				createdEntity.SetupGameEntity()
				createdEntity.AddComponents(map[string]ecs.Component{"position": ecs.PositionComponent{X: x, Y: y},
					"appearance": ecs.AppearanceComponent{Layer: 1, Character: "T", Color: "dark green", Name: "Troll"},
					"hitpoints": ecs.HitPointComponent{Hp: 20, MaxHP: 20},
					"block": ecs.BlockingComponent{},
					"movement": ecs.MovementComponent{},
					"basic_melee_ai": ecs.BasicMeleeAIComponent{},
					"attacker": ecs.AttackerComponent{Attack: 10, Defense: 7},
					"killable": ecs.KillableComponent{Name: "Remains of", Color: "dark red", Character: "%"}})
			} else if chance > 2 && chance <= 3 {
				// Create an Orc
				createdEntity = &ecs.GameEntity{}
				createdEntity.SetupGameEntity()
				createdEntity.AddComponents(map[string]ecs.Component{"position": ecs.PositionComponent{X: x, Y: y},
					"appearance": ecs.AppearanceComponent{Layer: 1, Character: "o", Color: "darker green", Name: "Orc"},
					"hitpoints": ecs.HitPointComponent{Hp: 15, MaxHP: 15},
					"block": ecs.BlockingComponent{},
					"movement": ecs.MovementComponent{},
					"basic_melee_ai": ecs.BasicMeleeAIComponent{},
					"attacker": ecs.AttackerComponent{Attack: 7, Defense: 5},
					"killable": ecs.KillableComponent{Name: "Remains of", Color: "dark red", Character: "%"}})
			} else if chance > 5 && chance <= 7 {
				// Create a Goblin
				createdEntity = &ecs.GameEntity{}
				createdEntity.SetupGameEntity()
				createdEntity.AddComponents(map[string]ecs.Component{"position": ecs.PositionComponent{X: x, Y: y},
					"appearance": ecs.AppearanceComponent{Layer: 1, Character: "g", Color: "green", Name: "Goblin"},
					"hitpoints": ecs.HitPointComponent{Hp: 5, MaxHP: 5},
					"block": ecs.BlockingComponent{},
					"movement": ecs.MovementComponent{},
					"basic_melee_ai": ecs.BasicMeleeAIComponent{},
					"attacker": ecs.AttackerComponent{Attack: 2, Defense: 2},
					"killable": ecs.KillableComponent{Name: "Remains of", Color: "dark red", Character: "%"}})
			} else if chance > 10 {
				// Create a reproducing Fungus
				createdEntity = &ecs.GameEntity{}
				createdEntity.SetupGameEntity()
				createdEntity.AddComponents(map[string]ecs.Component{"position": ecs.PositionComponent{X: x, Y: y},
					"appearance": ecs.AppearanceComponent{Layer: 1, Character: "f", Color: "yellow", Name: "Fungus"},
					"hitpoints": ecs.HitPointComponent{Hp: 5, MaxHP: 5},
					"block": ecs.BlockingComponent{},
					"reproducer": ecs.ReproducesComponent{MaxTimes: 8, TimesRemaining: 8, PercentChance: 25},
					"killable": ecs.KillableComponent{Name: "Remains of", Color: "yellow", Character: "."}})
			}

			entities = append(entities, createdEntity)
		} else {
			// No location was found after 50 tries, which means the map is quite full. Stop here and return.
			break
		}
	}

	return entities
}