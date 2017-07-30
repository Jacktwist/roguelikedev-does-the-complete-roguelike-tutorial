package ecs

import "math"

func GetBlockingEntitiesAtLocation(entities []*GameEntity, destinationX, destinationY int) *GameEntity {
	// Return any entities that are at the destination location which would block movement
	for _, e := range entities {
		if e.HasComponent("block") {
			if e.HasComponent("position") {
				positionComponent, _ := e.Components["position"].(PositionComponent)

				if positionComponent.X == destinationX && positionComponent.Y == destinationY {
					return e
				}
			}
		}
	}
	return nil
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
