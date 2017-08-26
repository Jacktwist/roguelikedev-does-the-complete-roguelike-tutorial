package ecs

import (
	"math"
)

func GetBlockingEntitiesAtLocation(entities []*GameEntity, destinationX, destinationY int) *GameEntity {
	// Return any entities that are at the destination location which would block movement
	for _, e := range entities {
		if e != nil {
			if e.HasComponent("block") {
				if e.HasComponent("position") {
					positionComponent, _ := e.Components["position"].(PositionComponent)

					if positionComponent.X == destinationX && positionComponent.Y == destinationY {
						return e
					}
				}
			}
		}
	}
	return nil
}

func GetEntitiesPresentAtLocation(entities []*GameEntity, x, y int) []*GameEntity {
	entitiesPresent := []*GameEntity{}

	for _, e := range entities {
		if e != nil {
			if e.HasComponents([]string{"position", "appearance"}) {
				pos, _ := e.Components["position"].(PositionComponent)

				if pos.X == x && pos.Y == y && !e.HasComponent("player") {
					// This entity is present at the currently examined location, so add its name to the list of present
					// entities
					entitiesPresent = append(entitiesPresent, e)
				}
			}
		}
	}

	return entitiesPresent
}

func GetEntityNamesPresentAtLocation(entities []*GameEntity, x, y int) string {
	entitiesPresent := []string{}

	for _, e := range entities {
		if e != nil {
			if e.HasComponents([]string{"position", "appearance"}) {
				pos, _ := e.Components["position"].(PositionComponent)
				appearance, _ := e.Components["appearance"].(AppearanceComponent)

				if pos.X == x && pos.Y == y {
					// This entity is present at the currently examined location, so add its name to the list of present
					// entities
					entitiesPresent = append(entitiesPresent, appearance.Name)
				}
			}
		}
	}

	entitiesList := ""
	if len(entitiesPresent) > 0 {
		listLen := len(entitiesPresent)
		for i := 0; i < len(entitiesPresent); i++ {
			if listLen == 1 || i == listLen-1 {
				entitiesList += entitiesPresent[i]
			} else {
				entitiesList += entitiesPresent[i] + ", "
			}

		}
	}
	return entitiesList
}

func getPlayerEntity(entities []*GameEntity) *GameEntity {
	// Searches through the game entity list, and returns the entity representing the player
	for _, e := range entities {
		if e.HasComponent("player") {
			return e
		}
	}
	return nil
}

func distanceTo(originX, originY, targetX, targetY int) int {
	dx := float64(targetX - originX)
	dy := float64(targetY - originY)

	result := int(Round(math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))))

	return result

}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

// Inventory related functions
func getExistingItemKey(entity *GameEntity, targetEntity *GameEntity) int{
	// Check to see if an entity has an identical item in their inventory. This is mostly useful for making sure
	// identical items get the same key assigned to them
	if entity.HasComponent("inventory") {
		inv, _ := entity.Components["inventory"].(InventoryComponent)

		if targetEntity.HasComponents([]string{"lootable", "appearance"}) {
			for i := 0; i < len(inv.Items); i++ {
				if inv.Items[i].HasComponent("appearance") {
					app, _ := inv.Items[i].Components["appearance"].(AppearanceComponent)
					targetApp, _ := targetEntity.Components["appearance"].(AppearanceComponent)

					if app.Name == targetApp.Name {
						lootable, _ := inv.Items[i].Components["lootable"].(LootableComponent)

						return lootable.Key
					}
				}
			}
		}
	}
	return 0
}

func FindItemWithKey(entity *GameEntity, keyCode int) *GameEntity {
	// Find an item in the entities inventory that has been assigned to the supplied key, and return it. If no such
	// entity exists, return nil.
	if entity.HasComponent("inventory") {
		inv, _ := entity.Components["inventory"].(InventoryComponent)

		for i := 0; i < len(inv.Items); i++ {
			if inv.Items[i].HasComponent("lootable") {
				lootable, _ := inv.Items[i].Components["lootable"].(LootableComponent)

				if lootable.Key == keyCode {
					return inv.Items[i]
				}
			}
		}
	}
	return nil
}

func CountItemInstances(entity, item *GameEntity) int {
	occurences := 0
	if entity.HasComponent("inventory") && item.HasComponent("appearance") {
		inv, _ := entity.Components["inventory"].(InventoryComponent)
		itemApp, _ := item.Components["appearance"].(AppearanceComponent)

		for i := 0; i < len(inv.Items); i++ {
			if inv.Items[i].HasComponent("appearance") {
				app, _ := inv.Items[i].Components["appearance"].(AppearanceComponent)

				if app.Name == itemApp.Name {
					occurences++
				}
			}
		}
	}
	return occurences
}