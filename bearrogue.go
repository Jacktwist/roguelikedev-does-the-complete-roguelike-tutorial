package main

import (
	blt "bearlibterminal"
	"bearrogue/camera"
	"bearrogue/ecs"
	"bearrogue/examinecursor"
	"bearrogue/fov"
	"bearrogue/gamemap"
	"bearrogue/ui"
	"fmt"
	"math/rand"
	"strconv"
)

const (
	WindowSizeX  = 100
	WindowSizeY  = 35
	ViewAreaX    = 75
	ViewAreaY    = 30
	MapWidth     = 100
	MapHeight    = 100
	Title        = "BearRogue"
	Font         = "fonts/UbuntuMono.ttf"
	FontSize     = 24
	PlayerTurn   = iota
	MobTurn      = iota
	MapLayer     = 0
	ActorLayer   = 2
	ItemLayer    = 3
	ExamineLayer = 4
)

var (
	version       string
	buildStamp    string
	gitHash       string
	player        *ecs.GameEntity
	entities      []*ecs.GameEntity
	gameMap       *gamemap.Map
	gameCamera    *camera.GameCamera
	fieldOfView   *fov.FieldOfVision
	gameTurn      int
	messageLog    ui.MessageLog
	examining     bool
	examineCursor *examinecursor.XCursor
	inMenu		  bool
	informationScreen bool
	inventoryKeys map[int]bool
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

	// Set the default menu state to false. This will be true when a menu is displayed, which will indicate that control
	// needs to be given to the menu, instead of the main game
	inMenu = false
	informationScreen = false

	inventoryKeys = map[int]bool{blt.TK_A: false,
		blt.TK_B: false,
		blt.TK_C: false,
		blt.TK_D: false,
		blt.TK_E: false,
		blt.TK_F: false,
		blt.TK_G: false,
		blt.TK_H: false,
		blt.TK_I: false,
		blt.TK_J: false,
		blt.TK_K: false,
		blt.TK_L: false,
		blt.TK_M: false,
		blt.TK_N: false,
		blt.TK_O: false,
		blt.TK_P: false,
		blt.TK_Q: false,
		blt.TK_R: false,
		blt.TK_S: false,
		blt.TK_T: false,
		blt.TK_U: false,
		blt.TK_V: false,
		blt.TK_W: false,
		blt.TK_X: false,
		blt.TK_Y: false,
		blt.TK_Z: false,
	}
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

		if key != blt.TK_CLOSE {
			if gameTurn == PlayerTurn {
				if player.HasComponents([]string{"movement", "controllable", "position"}) {
					handleInput(key, player)
				}
			}
		} else {
			break
		}

		if !inMenu {
			// Clear each Entity off the screen
			ecs.SystemClear(entities, gameCamera)

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
		} else {
			// Render the menu or information screen
			if inMenu && !informationScreen {
				renderinventory()
			}
		}
	}

	blt.Close()
}

func handleInput(key int, entity *ecs.GameEntity) {
	// Handle basic character movement in the four main directions, plus diagonals (and vim keys)

	var (
		dx, dy int
	)

	actionTaken := true

	if !inMenu {
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
			inventoryKeys = ecs.SystemPickupItem(player, entities, gameCamera, &messageLog, inventoryKeys)
		case blt.TK_I:
			inMenu = true
		case blt.TK_ESCAPE:
			// Cancel the current action and return the game state to normal
			actionTaken = false
			examining = false
			inMenu = false
			informationScreen = false
			if examineCursor != nil {
				examineCursor.Clear(gameCamera)
			}
			ui.ClearScreen(WindowSizeX, WindowSizeX)
		}
	} else {
		// If we are in a menu, bound controls will be different than standard.
		switch key {
		case blt.TK_ESCAPE:
			// Cancel the current action and return the game state to normal
			actionTaken = false
			examining = false
			if informationScreen {
				informationScreen = false
			} else {
				inMenu = false
			}

			if examineCursor != nil {
				examineCursor.Clear(gameCamera)
			}
			ui.ClearScreen(WindowSizeX, WindowSizeX)
		}

		selectedEntity := ecs.FindItemWithKey(player, key)

		if selectedEntity != nil {
			informationScreen = true
			renderInformationScreen(selectedEntity)
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
	if actionTaken && !examining && !inMenu {
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

func renderinventory() {
	ui.ClearScreen(WindowSizeX, WindowSizeX)

	if player.HasComponent("inventory") {
		inv, _ := player.Components["inventory"].(ecs.InventoryComponent)

		items := map[string]int{}

		// Get the counts of each type of item in the players inventory. If there is more than one of a type, and they
		// are stackable, group them together, to make the inventory more manageable
		// TODO: Refactor this to use the util function counting occurences of items.
		for i := 0; i < len(inv.Items); i++ {
			if inv.Items[i].HasComponents([]string{"appearance", "lootable"}) {
				app, _ := inv.Items[i].Components["appearance"].(ecs.AppearanceComponent)
				lootable, _ := inv.Items[i].Components["lootable"].(ecs.LootableComponent)

				key := string(ui.MapBltKeyCodesToRunes(lootable.Key))
				name := key + " - " + "[color=" + app.Color + "]" + app.Name + "[/color]"

				items[name]++
			}
		}
		ui.DisplayInventory(inv.Capacity, len(inv.Items), items)
	}
}

func renderInformationScreen(item *ecs.GameEntity) {
	ui.ClearScreen(WindowSizeX, WindowSizeX)

	if item.HasComponents([]string{"lootable", "appearance", "description"}) {
		app, _ := item.Components["appearance"].(ecs.AppearanceComponent)
		lootable, _ := item.Components["lootable"].(ecs.LootableComponent)
		desc, _ := item.Components["description"].(ecs.DescriptionComponent)

		key := string(ui.MapBltKeyCodesToRunes(lootable.Key))
		title := key + " - [color=" + app.Color + "]" + app.Name + "[/color]"

		occurences := ecs.CountItemInstances(player, item)

		ui.DisplayInformationScreen(title, desc.ShortDesc, desc.LongDesc, occurences, WindowSizeY)
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
			chance := rand.Intn(100)

			if chance >= 49 {
				// Create a healing potion
				createdEntity = &ecs.GameEntity{}
				createdEntity.SetupGameEntity()
				createdEntity.AddComponents(map[string]ecs.Component{"position": ecs.PositionComponent{X: x, Y: y},
					"appearance": ecs.AppearanceComponent{Layer: ItemLayer, Character: "!", Color: "dark red", Name: "Dark Red Potion"},
					"lootable":   ecs.LootableComponent{InInventory: false, ID: 1},
					"stackable":  ecs.StackableComponent{},
					"description": ecs.DescriptionComponent{ShortDesc:"An unmarked, single dose, vial of a dark red liquid."}})

				entities = append(entities, createdEntity)
			} else {
				// Create a healing potion
				createdEntity = &ecs.GameEntity{}
				createdEntity.SetupGameEntity()
				createdEntity.AddComponents(map[string]ecs.Component{"position": ecs.PositionComponent{X: x, Y: y},
					"appearance": ecs.AppearanceComponent{Layer: ItemLayer, Character: "!", Color: "light green", Name: "Bright Green Potion"},
					"lootable":   ecs.LootableComponent{InInventory: false, ID: 2},
					"stackable":  ecs.StackableComponent{},
					"description": ecs.DescriptionComponent{ShortDesc:"An unmarked, single dose, vial of a bright green liquid."}})

				entities = append(entities, createdEntity)
			}

		} else {
			// No location was found after 50 tries, which means the map is quite full. Stop here and return.
			break
		}
	}

	return entities
}
