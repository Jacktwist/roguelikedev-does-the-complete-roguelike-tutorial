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
	Name string
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

type HitPointComponent struct {
	Hp int
}

func (h HitPointComponent) IsComponent() bool {
	return true
}

// Attacker Component
type AttackerComponent struct {
	Attack int
	Defense int
}

func (a AttackerComponent) IsComponent() bool {
	return true
}

// Blocking Component
type BlockingComponent struct {

}

func (b BlockingComponent) IsComponent() bool {
	return true
}

// Random Movement Component - wanders aimlessly around the map
type RandomMovementComponent struct {

}

func (r RandomMovementComponent) IsComponent() bool {
	return true
}

