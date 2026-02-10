package components

import (
	"github.com/GiannisPettas/ember2D/internal/engine/entity"
)

// ComponentManager provides type-safe, generic storage for a single component type.
// Each component type gets its own manager.
//
// Usage:
//
//	positions := components.NewComponentManager[Position]()
//	positions.Add(playerID, Position{X: 100, Y: 50})
//	pos := positions.Get(playerID)  // returns *Position, no type assertion!
type ComponentManager[T any] struct {
	data map[entity.Entity]*T
}

// NewComponentManager creates a new ComponentManager for type T.
func NewComponentManager[T any]() *ComponentManager[T] {
	return &ComponentManager[T]{
		data: make(map[entity.Entity]*T),
	}
}

// Add attaches a component to an entity. Overwrites if already exists.
func (cm *ComponentManager[T]) Add(e entity.Entity, component T) {
	cm.data[e] = &component
}

// Get retrieves the component for an entity. Returns nil if not found.
func (cm *ComponentManager[T]) Get(e entity.Entity) *T {
	return cm.data[e]
}

// Remove detaches a component from an entity.
func (cm *ComponentManager[T]) Remove(e entity.Entity) {
	delete(cm.data, e)
}

// Has checks if an entity has this component.
func (cm *ComponentManager[T]) Has(e entity.Entity) bool {
	_, exists := cm.data[e]
	return exists
}

// Each iterates over all entities with this component.
func (cm *ComponentManager[T]) Each(fn func(entity.Entity, *T)) {
	for e, component := range cm.data {
		fn(e, component)
	}
}

// Count returns how many entities have this component.
func (cm *ComponentManager[T]) Count() int {
	return len(cm.data)
}
