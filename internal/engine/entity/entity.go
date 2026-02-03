package entity

import (
	"fmt"
)

type Entity struct {
	ID         string
	Tags       []string
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
		Tags:       make([]string, 0),
		Components: make(map[string]any),
	}

	// Add to world
	w.Entities[id] = entity

	return entity
}

// AddTag adds a tag to the entity (normalized to lowercase alphanumeric)
func (e *Entity) AddTag(tag string) {
	// Manually filter valid characters (faster than regex)
	tag = filterTag(tag)

	// Skip empty tags
	if tag == "" {
		return
	}

	if !e.HasTag(tag) {
		e.Tags = append(e.Tags, tag)
	}
}

// HasTag checks if the entity has a specific tag
func (e *Entity) HasTag(tag string) bool {
	// Early exit if no tags
	if len(e.Tags) == 0 {
		return false
	}

	// Normalize input
	tag = filterTag(tag)

	for _, t := range e.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// RemoveTag removes a tag from the entity
func (e *Entity) RemoveTag(tag string) {
	// Early exit if no tags
	if len(e.Tags) == 0 {
		return
	}

	// Normalize input
	tag = filterTag(tag)

	for i, t := range e.Tags {
		if t == tag {
			// Remove by swapping with last element and slicing
			e.Tags[i] = e.Tags[len(e.Tags)-1]
			e.Tags = e.Tags[:len(e.Tags)-1]
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
