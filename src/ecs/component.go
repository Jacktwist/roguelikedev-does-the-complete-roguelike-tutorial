package ecs

type Component interface {
	IsComponent() bool
}

// Player Component
type PlayerComponent struct {
}

func (pl PlayerComponent) IsComponent() bool {
	return true
}

// Position Component
type PositionComponent struct {
	X int
	Y int
}

func (pc PositionComponent) IsComponent() bool {
	return true
}

// Appearance Component
type AppearanceComponent struct {
	Color string
	Character string
	Layer int
}

func (a AppearanceComponent) IsComponent() bool {
	return true
}

// Movement Component
type MovementComponent struct {

}

func (m MovementComponent) IsComponent() bool {
	return true
}

// Controllable Component
type ControllableComponent struct {

}

func (c ControllableComponent) IsComponent() bool {
	return true
}

