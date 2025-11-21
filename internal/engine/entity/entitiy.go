package entity

type Entity struct {
	ID         string
	Components map[string]any
}

type World struct {
	Entities map[string]*Entity
}

func NewWorld() *World {
	return &World{
		Entities: make(map[string]*Entity),
	}
}

func (w *World) AddEntity(e *Entity) {
	w.Entities[e.ID] = e
}

func (w *World) GetEntity(id string) *Entity {
	return w.Entities[id]
}
