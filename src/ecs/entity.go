package ecs

import (
	"github.com/google/uuid"
	"fmt"
)

type GameEntity struct {
	gmUUID     uuid.UUID
	Components map[string]Component
}

func NewGameEntity() *GameEntity {
		gm := GameEntity {
			gmUUID: uuid.New(),
			Components: make(map[string]Component),
		}

	return &gm
}

func (e *GameEntity) HasComponent(componentName string) bool {
	if _, ok := e.Components[componentName]; ok {
		return true
	} else {
		return false
	}
}

func (e *GameEntity) HasComponents(componentNames []string) bool {
	containsAll := true
	for i := 0; i < len(componentNames); i++ {
		if !e.HasComponent(componentNames[i]) {
			containsAll = false
		}
	}
	return containsAll
}

func (e *GameEntity) AddComponent(name string, component Component) {
	e.Components[name] = component
}

func (e *GameEntity) RemoveComponent(componentName string) {
	_, ok := e.Components[componentName]

	if ok {
		delete(e.Components, componentName)
	}
}

func (e *GameEntity) GetComponent(componentName string) Component {
	if _, ok := e.Components[componentName]; ok {
		return e.Components[componentName]
	} else {
		return nil
}
}

func (e *GameEntity) Print() {
	fmt.Printf("uuid: %v\n", e.gmUUID)
	fmt.Printf("components: \n")
	for name, component := range e.Components {
			fmt.Printf("	%s: %+v\n", name, component)
	}
}

//func GetBlockingEntitiesAtLocation(entities []*GameEntity, destinationX, destinationY int) *GameEntity {
//	// Return any entities that are at the destination location which would block movement
//	for _, e := range entities {
//		if e.Blocks && e.X == destinationX && e.Y == destinationY {
//			return e
//		}
//	}
//	return nil
//}
