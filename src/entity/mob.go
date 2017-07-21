package entity

type Mob struct {
	Entity *GameEntity
}

func (m *Mob) Act() {
	// Allows the mob to take an action. Eventually, this will be replaced with an AI routine that will be unique for
	// the type of Mob this is, but for now it is a very simple action
}
