package entity

// Entity is just a numeric ID - lightweight and fast.
type Entity uint64

// World manages entity lifecycle: creation, destruction, and cleanup.
type World struct {
	nextID           Entity
	alive            map[Entity]bool
	entitiesToDelete []Entity
	tags             *TagManager
}

// NewWorld creates a new game world.
func NewWorld() *World {
	return &World{
		nextID:           0,
		alive:            make(map[Entity]bool),
		entitiesToDelete: make([]Entity, 0),
		tags:             NewTagManager(),
	}
}

// CreateEntity creates a new entity with a unique ID.
func (w *World) CreateEntity(tags ...string) Entity {
	id := w.nextID
	w.nextID++
	w.alive[id] = true
	for _, tag := range tags {
		w.tags.AddTag(id, tag)
	}
	return id
}

// DestroyEntity marks an entity for deletion (removed at end of frame).
func (w *World) DestroyEntity(e Entity) {
	if !w.alive[e] {
		return
	}
	w.alive[e] = false
	w.entitiesToDelete = append(w.entitiesToDelete, e)
}

// IsAlive checks if an entity exists and is not marked for deletion.
func (w *World) IsAlive(e Entity) bool {
	return w.alive[e]
}

// Cleanup removes all entities marked for deletion. Call at end of frame.
func (w *World) Cleanup() {
	for _, e := range w.entitiesToDelete {
		w.tags.RemoveAllTags(e)
		delete(w.alive, e)
	}
	w.entitiesToDelete = w.entitiesToDelete[:0]
}

// Tags returns the TagManager for this world.
func (w *World) Tags() *TagManager {
	return w.tags
}

// EntityCount returns the number of alive entities.
func (w *World) EntityCount() int {
	count := 0
	for _, isAlive := range w.alive {
		if isAlive {
			count++
		}
	}
	return count
}
