package main

import (
	blt "bearlibterminal"
	"bearrogue/camera"
	"bearrogue/ecs"
	"bearrogue/fov"
	"bearrogue/gamemap"
	"bearrogue/ui"
	"fmt"
	"math/rand"
	"strconv"
	"bearrogue/examinecursor"
)

const (
	WindowSizeX = 100
	WindowSizeY = 35
	ViewAreaX   = 75
	ViewAreaY   = 30
	MapWidth    = 100
	MapHeight   = 100
	Title       = "BearRogue"
	Font        = "fonts/UbuntuMono.ttf"
	FontSize    = 24
	PlayerTurn  = iota
	MobTurn     = iota
	MapLayer = 0
	ActorLayer = 2
	ItemLayer = 3
	ExamineLayer = 4
)

var (
	version		string
	buildStamp  string
	gitHash		string
	player      *ecs.GameEntity
	entities    []*ecs.GameEntity
	gameMap     *gamemap.Map
	gameCamera  *camera.GameCamera
	fieldOfView *fov.FieldOfVision
	gameTurn    int
	messageLog  ui.MessageLog
	examining 	bool
	examineCursor	*examinecursor.XCursor
)

func init() {
	blt.Open()

	fmt.Printf("BearRogue -- Version %s\n", version)
	fmt.Printf("Build Stamp: %s\n", buildStamp)
	fmt.Printf("Git Hash: %s\n", gitHash)

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
	player.AddComponent("inventory", ecs.InventoryComponent{Capacity: 32})

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
	}

	entities = append(entities, mapEntities...)

	// Set the current turn to the player, so they may act first
	gameTurn = PlayerTurn

	// Initialize a camera object
	gameCamera = &camera.GameCamera{X: 1, Y: 1, Width: ViewAreaX, Height: ViewAreaY}

	// Initialize a FoV object
	fieldOfView = &fov.FieldOfVision{}
	fieldOfView.Initialize()
	fieldOfView.SetTorchRadius(6)

	// Set up the messageLog, and output a "welcome" message
	messageLog = ui.MessageLog{MaxLength: 100}
	messageLog.InitMessages()

	// Set the default examining state to false, which means movement is normal. If this is true, only the examinecursor will
	// be moved, until the player cancels examine mode
	examining = false
}

func main() {
	// Main game loop

	renderSideBar()

	messageLog.SendMessage("You find yourself in the caverns of eternal sadness...you start to feel a little more sad.")
	renderMap()
	ecs.SystemRender(entities, gameCamera, gameMap)
	messageLog.PrintMessages(ViewAreaY, WindowSizeX, WindowSizeY)


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

		if examining {
			examineCursor.Draw(gameCamera)
		} else {
			// Only print messages if the player is not examining
			messageLog.PrintMessages(ViewAreaY, WindowSizeX, WindowSizeY)
		}
		renderSideBar()
	}

	blt.Close()
}

func handleInput(key int, entity *ecs.GameEntity) {
	// Handle basic character movement in the four main directions, plus diagonals (and vim keys)

	var (
		dx, dy int
	)

	actionTaken := true

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
	case blt.TK_X:
		// Look command - toggle the examining flag, and notify that this will not consume an action
		actionTaken = false
		examining = !examining

		if player.HasComponent("position") && examining {
			pos, _ := player.Components["position"].(ecs.PositionComponent)
			examineCursor = &examinecursor.XCursor{X: pos.X, Y: pos.Y, Character: "_", Layer: ExamineLayer}
		} else {
			examineCursor.Clear(gameCamera)
		}
	case blt.TK_COMMA:
		ecs.SystemPickupItem(player, entities, gameMap, &messageLog)
	case blt.TK_ESCAPE:
		// Cancel the current action and return the game state to normal
		actionTaken = false
		examining = false
		if examineCursor != nil {
			examineCursor.Clear(gameCamera)
		}
	}

	if examining {
		// Fire off examinecursor movement
		examine(dx, dy)
	} else {
		// Fire off the movement system
		ecs.SystemMovement(entity, dx, dy, entities, gameMap, &messageLog)
	}

	// Switch the game turn to the Mobs turn, if an action was taken. Some commands, like examine, or checking inventory
	// do not cost an action
	if actionTaken  && !examining {
		gameTurn = MobTurn
	}
}

func renderMap() {
	// Render the game map. If a tile is blocked and blocks sight, draw a '#', if it is not blocked, and does not block
	// sight, draw a '.'

	// First, set the every portion of the map seen by the camera to not visible. We'll decide what is visible based on
	// the torch radius. In the process, clear every camera visible Tile on the map as well
	for x := 0; x < MapWidth; x++ {
		for y := 0; y < MapHeight; y++ {
			gameMap.Tiles[x][y].Visible = false
		}
	}

	for x := 0; x < gameCamera.Width; x++ {
		for y := 0; y < gameCamera.Height; y++ {
			// Clear both our primary layers, so we don't get any strange artifacts from one layer or the other getting
			// cleared.
			for i := 0; i <= 2; i++ {
				blt.Layer(i)
				blt.Print(x, y, " ")
			}
		}
	}

	positionComponent, posOk := player.Components["position"].(ecs.PositionComponent)

	if posOk {
		gameCamera.MoveCamera(positionComponent.X, positionComponent.Y, MapWidth, MapHeight)

		// Next figure out what is visible to the player, and what is not.
		fieldOfView.RayCast(positionComponent.X, positionComponent.Y, gameMap)
	}

	// Now draw each tile that should appear on the screen, if its visible, or explored
	blt.Layer(MapLayer)
	for x := 0; x < gameCamera.Width; x++ {
		for y := 0; y < gameCamera.Height; y++ {
			mapX, mapY := gameCamera.X+x, gameCamera.Y+y

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

func renderSideBar() {
	blt.Layer(0)
	blt.ClearArea(ViewAreaX, 0, WindowSizeX, WindowSizeY)

	if player.HasComponents([]string{"appearance", "hitpoints"}) {
		playerAppearance, _ := player.Components["appearance"].(ecs.AppearanceComponent)
		playerHp, _ := player.Components["hitpoints"].(ecs.HitPointComponent)
		ui.PrintBasicCharacterInfo(playerAppearance.Name, ViewAreaX)
		ui.PrintStats(playerHp.Hp, playerHp.MaxHP, ViewAreaX)
	}
}

func examine(dx, dy int) {
	// Examine command - Creates a new cursor, that moves independently of the player, takes no actions, and will list
	// out any entities present at the location it is position on.
	examineCursor.Clear(gameCamera)

	examineCursor.Move(dx, dy, MapWidth, MapHeight, gameCamera)

	examineCursor.Draw(gameCamera)

	if gameMap.IsVisibleAndExplored(examineCursor.X, examineCursor.Y) {
		presentEntities := ecs.GetEntityNamesPresentAtLocation(entities, examineCursor.X, examineCursor.Y)
		if presentEntities != "" {
			ui.PrintToMessageArea(presentEntities, ViewAreaY, WindowSizeX, WindowSizeY, examineCursor.Layer)
		} else {
			tile := gameMap.Tiles[examineCursor.X][examineCursor.Y]
			if tile.IsWall() {
				ui.PrintToMessageArea("A cavern wall, made of some kind of rock", ViewAreaY, WindowSizeX, WindowSizeY, examineCursor.Layer)
			} else {
				ui.PrintToMessageArea("A cavern floor, covered in dirt and stones", ViewAreaY, WindowSizeX, WindowSizeY, examineCursor.Layer)
			}
		}
	} else {
		ui.PrintToMessageArea("You cannot see here...", ViewAreaY, WindowSizeX, WindowSizeY, examineCursor.Layer)
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
			if chance <= 5 {
				// Create a Troll
				createdEntity = &ecs.GameEntity{}
				createdEntity.SetupGameEntity()
				createdEntity.AddComponents(map[string]ecs.Component{"position": ecs.PositionComponent{X: x, Y: y},
					"appearance":     ecs.AppearanceComponent{Layer: ActorLayer, Character: "T", Color: "dark green", Name: "Troll"},
					"hitpoints":      ecs.HitPointComponent{Hp: 20, MaxHP: 20},
					"block":          ecs.BlockingComponent{},
					"movement":       ecs.MovementComponent{},
					"basic_melee_ai": ecs.BasicMeleeAIComponent{},
					"attacker":       ecs.AttackerComponent{Attack: 10, Defense: 7},
					"killable":       ecs.KillableComponent{Name: "Remains of", Color: "dark red", Character: "%"}})
			} else if chance > 5 && chance <= 20 {
				// Create an Orc
				createdEntity = &ecs.GameEntity{}
				createdEntity.SetupGameEntity()
				createdEntity.AddComponents(map[string]ecs.Component{"position": ecs.PositionComponent{X: x, Y: y},
					"appearance":     ecs.AppearanceComponent{Layer: ActorLayer, Character: "o", Color: "darker green", Name: "Orc"},
					"hitpoints":      ecs.HitPointComponent{Hp: 15, MaxHP: 15},
					"block":          ecs.BlockingComponent{},
					"movement":       ecs.MovementComponent{},
					"basic_melee_ai": ecs.BasicMeleeAIComponent{},
					"attacker":       ecs.AttackerComponent{Attack: 7, Defense: 5},
					"killable":       ecs.KillableComponent{Name: "Remains of", Color: "dark red", Character: "%"}})
			} else if chance > 20 && chance <= 70 {
				// Create a Goblin
				createdEntity = &ecs.GameEntity{}
				createdEntity.SetupGameEntity()
				createdEntity.AddComponents(map[string]ecs.Component{"position": ecs.PositionComponent{X: x, Y: y},
					"appearance":     ecs.AppearanceComponent{Layer: ActorLayer, Character: "g", Color: "green", Name: "Goblin"},
					"hitpoints":      ecs.HitPointComponent{Hp: 5, MaxHP: 5},
					"block":          ecs.BlockingComponent{},
					"movement":       ecs.MovementComponent{},
					"basic_melee_ai": ecs.BasicMeleeAIComponent{},
					"attacker":       ecs.AttackerComponent{Attack: 2, Defense: 2},
					"killable":       ecs.KillableComponent{Name: "Remains of", Color: "dark red", Character: "%"}})
			} else if chance > 70 {
				// Create a reproducing Fungus
				createdEntity = &ecs.GameEntity{}
				createdEntity.SetupGameEntity()
				createdEntity.AddComponents(map[string]ecs.Component{"position": ecs.PositionComponent{X: x, Y: y},
					"appearance": ecs.AppearanceComponent{Layer: ActorLayer, Character: "f", Color: "yellow", Name: "Fungus"},
					"hitpoints":  ecs.HitPointComponent{Hp: 5, MaxHP: 5},
					"block":      ecs.BlockingComponent{},
					"reproducer": ecs.ReproducesComponent{MaxTimes: 8, TimesRemaining: 8, PercentChance: 25},
					"killable":   ecs.KillableComponent{Name: "Remains of", Color: "yellow", Character: "."}})
			}

			entities = append(entities, createdEntity)
		} else {
			// No location was found after 50 tries, which means the map is quite full. Stop here and return.
			break
		}
	}

	// Next, populate some items (LootableComponent) in the dungeon in the same way. For now, just health potions will
	// be created.
	for i := 0; i < 15; i++ {
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
			// Create a healing potion
			createdEntity = &ecs.GameEntity{}
			createdEntity.SetupGameEntity()
			createdEntity.AddComponents(map[string]ecs.Component{"position": ecs.PositionComponent{X: x, Y: y},
				"appearance": ecs.AppearanceComponent{Layer: ItemLayer, Character: "!", Color: "dark red", Name: "Dark Red Potion"},
				"lootable": ecs.LootableComponent{InInventory: false}})

			entities = append(entities, createdEntity)
		} else {
			// No location was found after 50 tries, which means the map is quite full. Stop here and return.
			break
		}
	}

	return entities
}
