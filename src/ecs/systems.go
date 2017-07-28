package ecs

import (
	blt "bearlibterminal"
	"camera"
	"gamemap"
	"ui"
	"math/rand"
	"strconv"
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
	if entity.HasComponent("attacker") && entity != targetEntity{
		// Check to ensure the target entity has hitpoints. If it doesn't, check to see if it can be interacted with
		if targetEntity.HasComponents([]string{"hitpoints", "appearance"}) {

			eAppearanceComponent, _ := entity.Components["appearance"].(AppearanceComponent)
			tAppearanceComponent, _ := targetEntity.Components["appearance"].(AppearanceComponent)

			eAttackerComponent, _ := entity.Components["attacker"].(AttackerComponent)
			tAttackerComponent, _ := targetEntity.Components["attacker"].(AttackerComponent)

			// Simple attack algorithm (temporary): Attacking entitys attack value + d6 - defenders defense value
			attackModifier := rand.Intn(6)
			totalAttack := eAttackerComponent.Attack + attackModifier

			if totalAttack > tAttackerComponent.Defense {
				// The attack exceeded the defense of the target, so any excess should be applied as damage
				excess := totalAttack - tAttackerComponent.Defense

				tHitPointsComponent, _ := targetEntity.Components["hitpoints"].(HitPointComponent)

				tHitPointsComponent.Hp -= excess

				targetEntity.RemoveComponent("hitpoints")
				targetEntity.AddComponent("hitpoints", tHitPointsComponent)

				if entity.HasComponent("player") || targetEntity.HasComponent("player") {
					messageLog.SendMessage(eAppearanceComponent.Name + " attacks the " + tAppearanceComponent.Name + " for " + strconv.Itoa(excess) + " points of damage.")
				}

				// Check to see if this attack has reduced the targets HP to 0 or less
				if tHitPointsComponent.Hp <= 0 {
					// This entity has died, replace it with a corpse, and remove all movement and blocking components

					if entity.HasComponent("player") || targetEntity.HasComponent("player") {
						messageLog.SendMessage("The " + tAppearanceComponent.Name + " has been killed!")
					}

					tAppearanceComponent.Name = "Remains of " + tAppearanceComponent.Name
					tAppearanceComponent.Character = "%"
					tAppearanceComponent.Color = "dark red"
					tAppearanceComponent.Layer = 1

					targetEntity.RemoveComponent("appearance")
					targetEntity.AddComponent("appearance", tAppearanceComponent)

					targetEntity.RemoveComponents([]string{"movement", "attacker", "block", "random_movement", "hitpoints"})
				}
			}
		} else if targetEntity.HasComponent("appearance") {
			// The target cannot be attacked
			eAppearanceComponent, _ := entity.Components["appearance"].(AppearanceComponent)
			tAppearanceComponent, _ := targetEntity.Components["appearance"].(AppearanceComponent)

			if entity.HasComponent("player") || targetEntity.HasComponent("player") {
				messageLog.SendMessage(eAppearanceComponent.Name + " bumps into the " + tAppearanceComponent.Name + "\n")
			}
		}
	}
}
