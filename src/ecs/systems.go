package ecs

import (
	blt "bearlibterminal"
	"camera"
	"gamemap"
	"ui"
	"math/rand"
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

func SystemMovement(entity *GameEntity, dx, dy int, entities []*GameEntity, gameMap *gamemap.Map, messageLog *ui.MessageLog) {
	// Allow a moveable and controllable entity to move
	if entity.HasComponents([]string{"movement", "controllable", "position"}) {
		// If the current entity is controllable, moveable, and has a position, go ahead and move it
		positionComponent, _ := entity.Components["position"].(PositionComponent)

		if !gameMap.IsBlocked(positionComponent.X + dx, positionComponent.Y + dy) {
			target := GetBlockingEntitiesAtLocation(entities, positionComponent.X+dx, positionComponent.Y+dy)
			if target != nil {
				SystemAttack(entity, target, messageLog)
			} else {
				positionComponent.X += dx
				positionComponent.Y += dy

				entity.RemoveComponent("position")
				entity.AddComponent("position", positionComponent)
			}
		}
	} else {
		// Otherwise, just give it random movement for now
		SystemRandomMovement(entity, entities, gameMap, messageLog)
	}
}

func SystemRandomMovement(entity *GameEntity, entities []*GameEntity, gameMap *gamemap.Map, messageLog *ui.MessageLog) {
	if entity.HasComponents([]string{"movement", "position", "random_movement"}) {
		positionComponent, _ := entity.Components["position"].(PositionComponent)

		// Choose a random (x, y) such that -1 <= x <= 1 and -1 <= y <= 1
		dx := rand.Intn(3) + -1
		dy := rand.Intn(3) + -1

		if !gameMap.IsBlocked(positionComponent.X + dx, positionComponent.Y + dy) {
			target := GetBlockingEntitiesAtLocation(entities, positionComponent.X+dx, positionComponent.Y+dy)
			if target != nil {
				SystemAttack(entity, target, messageLog)
			} else {
				positionComponent.X += dx
				positionComponent.Y += dy

				entity.RemoveComponent("position")
				entity.AddComponent("position", positionComponent)
			}
		}
	}
}

func SystemAttack(entity *GameEntity, targetEntity *GameEntity, messageLog *ui.MessageLog) {
	// Initiate an attack against another entity
	if entity.HasComponent("attacker") {
		// Check to ensure the target entity has hitpoints. If it doesn't, check to see if it can be interacted with
		if targetEntity.HasComponents([]string{"hitpoints", "appearance"}) {
			appearanceComponent, _ := targetEntity.Components["appearance"].(AppearanceComponent)
			messageLog.SendMessage("You kick the " + appearanceComponent.Name + " in the shins.")
		} else if targetEntity.HasComponent("appearance") {
			// The target cannot be attacked
			appearanceComponent, _ := targetEntity.Components["appearance"].(AppearanceComponent)

			messageLog.SendMessage("You bump into the " + appearanceComponent.Name + "\n")
		}
	}
}
