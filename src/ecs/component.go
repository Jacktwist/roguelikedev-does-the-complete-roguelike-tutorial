package ecs

type Component interface {
	IsComponent() bool
	IsAIComponent() bool
}

// Player Component
type PlayerComponent struct {
}

func (pl PlayerComponent) IsComponent() bool {
	return true
}

func (pl PlayerComponent) IsAIComponent() bool {
	return false
}

// Position Component
type PositionComponent struct {
	X int
	Y int
}

func (pc PositionComponent) IsComponent() bool {
	return true
}

func (pc PositionComponent) IsAIComponent() bool {
	return false
}

// Appearance Component
type AppearanceComponent struct {
	Color string
	Character string
	Layer int
	Name string
}

func (a AppearanceComponent) IsComponent() bool {
	return true
}

func (a AppearanceComponent) IsAIComponent() bool {
	return false
}

// Movement Component
type MovementComponent struct {

}

func (m MovementComponent) IsComponent() bool {
	return true
}

func (m MovementComponent) IsAIComponent() bool {
	return false
}

// Controllable Component
type ControllableComponent struct {

}

func (c ControllableComponent) IsComponent() bool {
	return true
}

func (c ControllableComponent) IsAIComponent() bool {
	return false
}

type HitPointComponent struct {
	Hp int
	MaxHP int
}

func (h HitPointComponent) IsComponent() bool {
	return true
}

func (h HitPointComponent) IsAIComponent() bool {
	return false
}

// Attacker Component
type AttackerComponent struct {
	Attack int
	Defense int
}

func (a AttackerComponent) IsComponent() bool {
	return true
}

func (a AttackerComponent) IsAIComponent() bool {
	return false
}

// Blocking Component
type BlockingComponent struct {

}

func (b BlockingComponent) IsComponent() bool {
	return true
}

func (b BlockingComponent) IsAIComponent() bool {
	return false
}

// Random Movement Component - wanders aimlessly around the map
type RandomMovementComponent struct {

}

func (r RandomMovementComponent) IsComponent() bool {
	return true
}

func (r RandomMovementComponent) IsAIComponent() bool {
	return true
}

// Basic Melee Attack AI Component
type BasicMeleeAIComponent struct {
	target *GameEntity
}

func (b BasicMeleeAIComponent) IsComponent() bool {
	return true
}

func (b BasicMeleeAIComponent) IsAIComponent() bool {
	return true
}



