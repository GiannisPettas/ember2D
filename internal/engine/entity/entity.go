package entity

import "fmt"

type Entity struct {
	ID         string
	Components map[string]any
}

type World struct {
	Entities  map[string]*Entity
	idCounter int //counter for generating IDs
}

func NewWorld() *World {
	return &World{
		Entities:  make(map[string]*Entity),
		idCounter: 0, //start from 0
	}
}

func (w *World) GetEntity(id string) *Entity {
	return w.Entities[id]
}

// Factory method for creating entities with auto-generated IDs
func (w *World) CreateEntity(prefix string) *Entity {
	id := fmt.Sprintf("%s_%d", prefix, w.idCounter)
	w.idCounter++

	// Create entity
	entity := &Entity{
		ID:         id,
		Components: make(map[string]any),
	}

	// Add to world
	w.Entities[id] = entity

	return entity
}
