package main

import (
	blt "bearlibterminal"
	"camera"
	"ecs"
	"fov"
	"gamemap"
	"strconv"
	"ui"
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
	player = ecs.NewGameEntity()
	player.AddComponent("player", ecs.PlayerComponent{})
	player.AddComponent("position", ecs.PositionComponent{X: 0, Y: 0})
	player.AddComponent("appearance", ecs.AppearanceComponent{Color: "white", Character: "@", Layer: 1})
	player.AddComponent("movement", ecs.MovementComponent{})
	player.AddComponent("controllable", ecs.ControllableComponent{})

	entities = append(entities, player)

	// Create a GameMap, and initialize it (and set the player position within it, for now)
	gameMap = &gamemap.Map{Width: MapWidth, Height: MapHeight}
	gameMap.InitializeMap()

	playerX, playerY:= gameMap.GenerateCavern()

	if player.HasComponent("position") {
		positionComponent, _ := player.Components["position"].(ecs.PositionComponent)
		positionComponent.X = playerX
		positionComponent.Y = playerY
		player.RemoveComponent("position")
		player.AddComponent("position", positionComponent)
		player.Print()
	}

	//entities = append(entities, mapEntities...)

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

		//if gameTurn == MobTurn {
		//	for _, e := range entities {
		//		if e != player {
		//			if gameMap.Tiles[e.X][e.Y].Visible {
		//				// Check to ensure that the ecs is visible before allowing it to message the player
		//				// This will change soon, as entities will act whether the player can see them or not.
		//				messageLog.SendMessage("The " + e.Name + " waits patiently.")
		//			}
		//		}
		//	}
		//	gameTurn = PlayerTurn
		//}

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

	// Check to ensure that the tile the player is trying to move in to is a valid move (not blocked)
	positionComponent, _ := player.Components["position"].(ecs.PositionComponent)

	if !gameMap.IsBlocked(positionComponent.X + dx, positionComponent.Y + dy) {
		//target := ecs.GetBlockingEntitiesAtLocation(entities, player.X + dx, player.Y + dy)
		//if target != nil {
		//	messageLog.SendMessage("You harmlessly bump into the " + target.Name)
		//} else {
		//	player.Move(dx, dy)
		//}
		ecs.SystemMovement(entity, dx, dy)
	}

	// Switch the game turn to the Mobs turn
	//gameTurn = MobTurn
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

