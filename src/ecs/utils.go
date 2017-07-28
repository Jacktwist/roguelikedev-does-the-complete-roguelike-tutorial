package ecs

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
