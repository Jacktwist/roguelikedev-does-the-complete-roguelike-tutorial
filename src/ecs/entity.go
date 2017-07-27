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
	// Check to see if the entity has the given component
	if _, ok := e.Components[componentName]; ok {
		return true
	} else {
		return false
	}
}

func (e *GameEntity) HasComponents(componentNames []string) bool {
	// Check to see if the entity has the given components
	containsAll := true
	for i := 0; i < len(componentNames); i++ {
		if !e.HasComponent(componentNames[i]) {
			containsAll = false
		}
	}
	return containsAll
}

func (e *GameEntity) AddComponent(name string, component Component) {
	// Add a single component to the entity
	e.Components[name] = component
}

func (e *GameEntity) AddComponents(components map[string]Component) {
	// Add several (or one) components to the entity
	for name, component := range components {
		e.Components[name] = component
	}
}

func (e *GameEntity) RemoveComponent(componentName string) {
	// Remove of a component from the entity
	_, ok := e.Components[componentName]

	if ok {
		delete(e.Components, componentName)
	}
}

func (e *GameEntity) GetComponent(componentName string) Component {
	// Return the named component from the entity, if present
	if _, ok := e.Components[componentName]; ok {
		return e.Components[componentName]
	} else {
		return nil
}
}

func (e *GameEntity) Print() {
	// Print a string representation of the entity, including all components
	fmt.Printf("uuid: %v\n", e.gmUUID)
	fmt.Printf("components: \n")
	for name, component := range e.Components {
			fmt.Printf("	%s: %+v\n", name, component)
	}
}
