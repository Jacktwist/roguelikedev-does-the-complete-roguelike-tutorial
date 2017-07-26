package ecs

import (
	blt "bearlibterminal"
	"camera"
	"gamemap"
)

func SystemRender(entities []*GameEntity, camera *camera.GameCamera, gameMap *gamemap.Map) {
	// Render all renderable entities to the screen
	for _, e := range entities {
		if e.HasComponents([]string{"position", "appearance"}) {
			pos, _ := e.Components["position"].(PositionComponent)
			app, _ := e.Components["appearance"].(AppearanceComponent)

			cameraX, cameraY := camera.ToCameraCoordinates(pos.X, pos.Y)

			if gameMap.Tiles[pos.X][pos.Y].Visible{
				blt.Layer(app.Layer)
				blt.Color(blt.ColorFromName(app.Color))
				blt.Print(cameraX, cameraY, app.Character)
			}
		}
	}
}

func SystemClear(entities []*GameEntity, camera *camera.GameCamera) {
	for _, e := range entities {
		if e.HasComponents([]string{"position", "appearance"}) {
			// Clear the entity from the screen. This only applies to entities that have a position and an
			// appearance
			positionComponent, _ := e.Components["position"].(PositionComponent)
			appearanceComponent, _ := e.Components["appearance"].(AppearanceComponent)

			mapX, mapY := camera.ToCameraCoordinates(positionComponent.X, positionComponent.Y)

			blt.Layer(appearanceComponent.Layer)
			blt.Print(mapX, mapY, " ")
		}
	}
}

func SystemMovement(entity *GameEntity, dx, dy int) {
	// Allow a moveable and controllable entity to move
	if entity.HasComponents([]string{"movement", "controllable", "position"}) {
		positionComponent, _ := entity.Components["position"].(PositionComponent)

		positionComponent.X += dx
		positionComponent.Y += dy

		entity.RemoveComponent("position")
		entity.AddComponent("position", positionComponent)
	}
}
