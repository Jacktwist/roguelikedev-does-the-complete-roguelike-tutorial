package ecs

import (
	blt "bearlibterminal"
	"bearrogue/camera"
	"bearrogue/gamemap"
	"bearrogue/ui"
	"math/rand"
	"strconv"
)

const (
	CorpseLayer = 2
)

func SystemRender(entities []*GameEntity, camera *camera.GameCamera, gameMap *gamemap.Map) {
	// Render all renderable entities to the screen
	for _, e := range entities {
		if e != nil {
			if e.HasComponents([]string{"position", "appearance"}) {
				pos, _ := e.Components["position"].(PositionComponent)
				app, _ := e.Components["appearance"].(AppearanceComponent)

				SystemClearAt(e, camera, pos.X, pos.Y)

				cameraX, cameraY := camera.ToCameraCoordinates(pos.X, pos.Y)

				if gameMap.Tiles[pos.X][pos.Y].Visible {
					blt.Layer(app.Layer)
					blt.Color(blt.ColorFromName(app.Color))
					blt.Print(cameraX, cameraY, app.Character)
				}
			}
		}
	}
}

func SystemClear(entities []*GameEntity, camera *camera.GameCamera) {
	for _, e := range entities {
		if e.HasComponents([]string{"position","appearance"}) {
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

func SystemClearAt(entity *GameEntity, camera *camera.GameCamera, x, y int) {
	// Clear an entity that may not have a position any longer
	if entity.HasComponent("appearance") {

		appearanceComponent, _ := entity.Components["appearance"].(AppearanceComponent)

		layer := blt.TK_LAYER
		blt.Layer(appearanceComponent.Layer)
		cameraX, cameraY := camera.ToCameraCoordinates(x, y)
		blt.Print(cameraX, cameraY, " ")
		blt.Layer(layer)
	}
}

func SystemMovement(entity *GameEntity, dx, dy int, entities []*GameEntity, gameMap *gamemap.Map, messageLog *ui.MessageLog) {
	// Allow a moveable and controllable entity to move
	if entity.HasComponents([]string{"movement", "controllable", "position"}) {
		// If the current entity is controllable, moveable, and has a position, go ahead and move it
		positionComponent, _ := entity.Components["position"].(PositionComponent)

		if !gameMap.IsBlocked(positionComponent.X+dx, positionComponent.Y+dy) {
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
		// Check if the entity has an AI component. If it does, use that for movement
		aiComponent := entity.HasAIComponent()
		if aiComponent != "" {
			switch aiComponent {
			case "random_movement":
				SystemRandomMovement(entity, entities, gameMap, messageLog)
			case "basic_melee_ai":
				SystemBasicMeleeAI(entity, entities, gameMap, messageLog)
			}
		}
	}
}

func SystemRandomMovement(entity *GameEntity, entities []*GameEntity, gameMap *gamemap.Map, messageLog *ui.MessageLog) {
	if entity.HasComponents([]string{"movement", "position"}) {

		positionComponent, _ := entity.Components["position"].(PositionComponent)

		// Choose a random (x, y) such that -1 <= x <= 1 and -1 <= y <= 1
		dx := rand.Intn(3) + -1
		dy := rand.Intn(3) + -1

		if !gameMap.IsBlocked(positionComponent.X+dx, positionComponent.Y+dy) {
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

func SystemBasicMeleeAI(entity *GameEntity, entities []*GameEntity, gameMap *gamemap.Map, messageLog *ui.MessageLog) {
	// This is the most basic AI available. The entity will choose a target, and move towards that target until it is
	// right next to it, then it will repeatedly attack the target. It chooses the closest viable target for its attacks
	if entity.HasComponents([]string{"position", "movement", "appearance", "basic_melee_ai"}) {
		//First, check to ensure the entity is within the players line of sight
		positionComponent, _ := entity.Components["position"].(PositionComponent)
		appearanceComponent, _ := entity.Components["appearance"].(AppearanceComponent)

		if gameMap.IsVisibleToPlayer(positionComponent.X, positionComponent.Y) {
			// The entity is currently within the players field of vision, it should do something
			// First, pick a target (this will usually be the player, but maybe not always)
			basicMeleeAi, _ := entity.Components["basic_melee_ai"].(BasicMeleeAIComponent)

			// For now, use the player
			target := getPlayerEntity(entities)

			targetPositionComponent, _ := target.Components["position"].(PositionComponent)

			oldTarget := basicMeleeAi.target

			// Set the target
			basicMeleeAi.target = target

			if oldTarget != basicMeleeAi.target {
				targetAppearanceComponent, _ := basicMeleeAi.target.Components["appearance"].(AppearanceComponent)
				messageLog.SendMessage("The [color=" + appearanceComponent.Color + "]" + appearanceComponent.Name + "[/color] throws an angry glare at [color=" + targetAppearanceComponent.Color + "]" + targetAppearanceComponent.Name + "[/color]!")
			}

			entity.RemoveComponent("basic_melee_ai")
			entity.AddComponent("basic_melee_ai", basicMeleeAi)

			// Now that the entity has a target, move towards it
			distance := distanceTo(positionComponent.X, positionComponent.Y, targetPositionComponent.X, targetPositionComponent.Y)

			dx := int(Round((float64(targetPositionComponent.X) - float64(positionComponent.X)) / float64(distance)))
			dy := int(Round((float64(targetPositionComponent.Y) - float64(positionComponent.Y)) / float64(distance)))

			if !gameMap.IsBlocked(positionComponent.X+dx, positionComponent.Y+dy) {
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
			// The entity is not currently visible to the player, so it should just shuffle around randomly for now
			SystemRandomMovement(entity, entities, gameMap, messageLog)
		}
	}

}

func SystemAttack(entity *GameEntity, targetEntity *GameEntity, messageLog *ui.MessageLog) {
	// Initiate an attack against another entity
	if entity.HasComponent("attacker") && entity != targetEntity {
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
					messageLog.SendMessage("[color=" + eAppearanceComponent.Color + "]" + eAppearanceComponent.Name + "[/color] attacks the [color=" + tAppearanceComponent.Color + "]" + tAppearanceComponent.Name + "[/color] for " + strconv.Itoa(excess) + " points of damage.")
				}

				// Check to see if this attack has reduced the targets HP to 0 or less
				if tHitPointsComponent.Hp <= 0 {
					// This entity has died, replace it with a corpse, and remove all movement and blocking components
					if targetEntity.HasComponent("killable") {
						if entity.HasComponent("player") || targetEntity.HasComponent("player") {
							messageLog.SendMessage("The [color=" + tAppearanceComponent.Color + "]" + tAppearanceComponent.Name + "[/color] has been killed!")
						}

						killableComponent, _ := targetEntity.Components["killable"].(KillableComponent)

						tAppearanceComponent.Name = killableComponent.Name + " " + tAppearanceComponent.Name
						tAppearanceComponent.Character = killableComponent.Character
						tAppearanceComponent.Color = killableComponent.Color
						tAppearanceComponent.Layer = CorpseLayer

						targetEntity.RemoveComponent("appearance")
						targetEntity.AddComponent("appearance", tAppearanceComponent)

						targetEntity.RemoveComponents([]string{"movement", "attacker", "block", "random_movement", "hitpoints", "reproducer"})
					}
				}
			} else {
				if entity.HasComponent("player") || targetEntity.HasComponent("player") {
					messageLog.SendMessage("[color=" + eAppearanceComponent.Color + "]" + eAppearanceComponent.Name + "[/color] attacks the [color=" + tAppearanceComponent.Color + "]" + tAppearanceComponent.Name + "[/color], but does no damage!")
				}
			}
		} else if targetEntity.HasComponent("appearance") {
			// The target cannot be attacked
			eAppearanceComponent, _ := entity.Components["appearance"].(AppearanceComponent)
			tAppearanceComponent, _ := targetEntity.Components["appearance"].(AppearanceComponent)

			if entity.HasComponent("player") || targetEntity.HasComponent("player") {
				messageLog.SendMessage("[color=" + eAppearanceComponent.Color + "]" + eAppearanceComponent.Name + "[/color] bumps into the [color=" + tAppearanceComponent.Color + "]" + tAppearanceComponent.Name + "[/color]\n")
			}
		}
	}
}

func SystemReproduce(entity *GameEntity, entities []*GameEntity, gameMap *gamemap.Map, messageLog *ui.MessageLog) *GameEntity {
	if entity.HasComponent("reproducer") {
		reproducerComponent, _ := entity.Components["reproducer"].(ReproducesComponent)

		chance := rand.Intn(100)

		if reproducerComponent.TimesRemaining > 0 && chance <= reproducerComponent.PercentChance {
			// This entity can still reproduce, so do so

			positionComponent, _ := entity.Components["position"].(PositionComponent)

			// Randomly generate a direction to reproduce in
			x := (rand.Intn(3) + -1) + positionComponent.X
			y := (rand.Intn(3) + -1) + positionComponent.Y

			if !gameMap.IsBlocked(x, y) {
				target := GetBlockingEntitiesAtLocation(entities, x, y)
				if target == nil {
					// There is nothing blocking the new entity, so go ahead and create it
					createdEntity := &GameEntity{}
					createdEntity.SetupGameEntity()
					//createdEntity.Components = entity.Components

					for name, e := range entity.Components {
						createdEntity.AddComponent(name, e)
					}

					// Update the position and number of reproductions
					rPositionComponent, _ := createdEntity.Components["position"].(PositionComponent)
					rReproducerComponent, _ := createdEntity.Components["reproducer"].(ReproducesComponent)

					rPositionComponent.X = x
					rPositionComponent.Y = y
					rReproducerComponent.TimesRemaining = rReproducerComponent.TimesRemaining - 2
					rReproducerComponent.PercentChance = int(reproducerComponent.PercentChance / 2)

					createdEntity.RemoveComponents([]string{"position", "reproducer"})

					createdEntity.AddComponents(map[string]Component{"position": rPositionComponent, "reproducer": rReproducerComponent})

					reproducerComponent.TimesRemaining -= 1
					entity.RemoveComponent("reproducer")
					entity.AddComponent("reproducer", reproducerComponent)

					return createdEntity

				}
			}

			reproducerComponent.TimesRemaining -= 1
			entity.RemoveComponent("reproducer")
			entity.AddComponent("reproducer", reproducerComponent)
		}
	}
	return nil
}

func SystemPickupItem(entity *GameEntity, entities []*GameEntity, camera *camera.GameCamera, messageLog *ui.MessageLog, inventoryKeys map[int]bool) map[int]bool {
	if entity.HasComponents([]string{"inventory", "position", "appearance"}) {
		inv, _ := entity.Components["inventory"].(InventoryComponent)
		pos, _ := entity.Components["position"].(PositionComponent)
		app, _ := entity.Components["appearance"].(AppearanceComponent)

		entitiesPresent := GetEntitiesPresentAtLocation(entities, pos.X, pos.Y)

		if len(entitiesPresent) > 0 {
			// For now, this assumes one entity per tile, which will obviously need to change
			targetEntity := entitiesPresent[0]

			if targetEntity.HasComponents([]string{"appearance", "position"}) {
				targetPosition, _ := targetEntity.Components["position"].(PositionComponent)
				targetAppearance, _ := targetEntity.Components["appearance"].(AppearanceComponent)

				if targetEntity.HasComponent("lootable") {
					targetLootable, _ := targetEntity.Components["lootable"].(LootableComponent)

					// Make sure the lootable is not currently in an inventory
					if len(inv.Items) < inv.Capacity && !targetLootable.InInventory {
						// Transfer the lootable entity to the players inventory
						targetLootable.InInventory = true
						targetLootable.Owner = entity

						key := getExistingItemKey(entity, targetEntity)

						if key != 0 {
							targetLootable.Key = key
						} else {
							// There was no existing identical item in the inventory, so we need to assign a key to this
							// one. Pull the key from a pool of possible, non-assigned keys
							for k, v := range inventoryKeys {
								if !v {
									targetLootable.Key = k
									break
								}
							}
							// Make sure we mark the key used as not available
							inventoryKeys[targetLootable.Key] = true
						}

						targetEntity.RemoveComponents([]string{"lootable", "position"})
						targetEntity.AddComponent("lootable", targetLootable)

						inv.Items = append(inv.Items, targetEntity)

						entity.RemoveComponent("inventory")
						entity.AddComponent("inventory", inv)

						SystemClearAt(targetEntity, camera, targetPosition.X, targetPosition.Y)

						messageLog.SendMessage(app.Name + " picks up the [color=" + targetAppearance.Color + "]" + targetAppearance.Name + "[/color]")
					} else {
						if entity.HasComponent("player") {
							messageLog.SendMessage("Your inventory is full, and you cannot pick up the ")
						}
					}
				} else {
					// The entity present is not lootable, notify, and do not add it to inventory
					messageLog.SendMessage("Cannot pick up that [color=" + targetAppearance.Color + "]" + targetAppearance.Name + "[/color]")
				}
			}
		} else {
			messageLog.SendMessage("There is nothing to pick up here!")
		}

	}
	return inventoryKeys
}
