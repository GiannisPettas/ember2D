package entity

import (
	"fmt"
)

// Entity represents an object in the world.
type Entity struct {
	ID         string
	tags       []string
	Components map[string]any
	isAlive    bool
}

// World manages entities and their lifecycle.
type World struct {
	Entities         map[string]*Entity
	tagIndex         map[string]map[string]*Entity // tag -> entityID -> *Entity
	idCounter        int                           //counter for generating IDs
	entitiesToDelete []string
}

func NewWorld() *World {
	return &World{
		Entities:         make(map[string]*Entity),
		tagIndex:         make(map[string]map[string]*Entity),
		idCounter:        0, //start from 0
		entitiesToDelete: make([]string, 0),
	}
}

func (w *World) GetEntity(id string) *Entity {
	return w.Entities[id]
}

// AddTag adds a tag to an entity and updates the index
func (w *World) AddTag(entity *Entity, tag string) {
	tag = filterTag(tag)
	if tag == "" {
		return
	}
	// Only add if entity doesn't already have this tag
	if entity.hasTagInternal(tag) {
		return
	}
	// Add to entity's internal tags
	entity.addTagInternal(tag)
	// Update index
	if w.tagIndex[tag] == nil {
		w.tagIndex[tag] = make(map[string]*Entity)
	}
	w.tagIndex[tag][entity.ID] = entity
}

// RemoveTag removes a tag from an entity and updates the index
func (w *World) RemoveTag(entity *Entity, tag string) {
	tag = filterTag(tag)
	if tag == "" {
		return
	}
	entity.removeTagInternal(tag)
	// Update index
	if w.tagIndex[tag] != nil {
		delete(w.tagIndex[tag], entity.ID)
	}
}

// HasTag checks if an entity has a specific tag
func (w *World) HasTag(entity *Entity, tag string) bool {
	tag = filterTag(tag)
	return entity.hasTagInternal(tag)
}

// GetTags returns a copy of the entity's tags
func (w *World) GetTags(entity *Entity) []string {
	result := make([]string, len(entity.tags))
	copy(result, entity.tags)
	return result
}

// GetEntitiesByTag returns all entities with the specified tag (O(1) lookup)
func (w *World) GetEntitiesByTag(tag string) []*Entity {
	tag = filterTag(tag)
	entityMap := w.tagIndex[tag]
	if entityMap == nil {
		return nil
	}
	result := make([]*Entity, 0, len(entityMap))
	for _, entity := range entityMap {
		result = append(result, entity)
	}
	return result
}

// DestroyEntity marks an entity for deletion
func (w *World) DestroyEntity(id string) {
	entity := w.Entities[id]
	if entity == nil || !entity.isAlive {
		return
	}
	entity.isAlive = false
	w.entitiesToDelete = append(w.entitiesToDelete, id)
}

// IsAlive checks if entity is not marked for deletion
func (w *World) IsAlive(entity *Entity) bool {
	return entity != nil && entity.isAlive
}

// Cleanup removes marked entities (call at end of frame)
func (w *World) Cleanup() {
	for _, id := range w.entitiesToDelete {
		entity := w.Entities[id]
		if entity == nil {
			continue
		}
		for _, tag := range entity.tags {
			if w.tagIndex[tag] != nil {
				delete(w.tagIndex[tag], entity.ID)
			}
		}
		delete(w.Entities, id)
	}
	w.entitiesToDelete = w.entitiesToDelete[:0]
}

// Factory method for creating entities with auto-generated IDs
func (w *World) CreateEntity(prefix string) *Entity {
	id := fmt.Sprintf("%s_%d", prefix, w.idCounter)
	w.idCounter++

	// Create entity
	entity := &Entity{
		ID:         id,
		tags:       make([]string, 0),
		Components: make(map[string]any),
		isAlive:    true,
	}

	// Add to world
	w.Entities[id] = entity
	// Auto-add prefix as tag for easy querying
	w.AddTag(entity, prefix) // Use World method

	return entity
}

// AddTag adds a tag to the entity (normalized to lowercase alphanumeric)
func (e *Entity) addTagInternal(tag string) {
	e.tags = append(e.tags, tag)
}

// HasTag checks if the entity has a specific tag
func (e *Entity) hasTagInternal(tag string) bool {
	// Early exit if no tags
	if len(e.tags) == 0 {
		return false
	}

	for _, t := range e.tags {
		if t == tag {
			return true
		}
	}
	return false
}

// RemoveTag removes a tag from the entity
func (e *Entity) removeTagInternal(tag string) {
	// Early exit if no tags
	if len(e.tags) == 0 {
		return
	}

	for i, t := range e.tags {
		if t == tag {
			// Remove by swapping with last element and slicing
			e.tags[i] = e.tags[len(e.tags)-1]
			e.tags = e.tags[:len(e.tags)-1]
			return
		}
	}
}

// filterTag normalizes and filters: lowercase, keeps only a-z, 0-9, underscore
func filterTag(s string) string {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]

		// Convert uppercase to lowercase
		if c >= 'A' && c <= 'Z' {
			c = c + 32 // 'A'(65) + 32 = 'a'(97)
		}

		// Keep only valid characters
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' {
			result = append(result, c)
		}
	}
	return string(result)
}
