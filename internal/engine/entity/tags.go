package entity

// TagManager handles entity tags with O(1) lookup via index.
type TagManager struct {
	entityTags map[Entity]map[string]bool // entity -> set of tags
	tagIndex   map[string]map[Entity]bool // tag -> set of entities (reverse index)
}

// NewTagManager creates a new TagManager.
func NewTagManager() *TagManager {
	return &TagManager{
		entityTags: make(map[Entity]map[string]bool),
		tagIndex:   make(map[string]map[Entity]bool),
	}
}

// AddTag adds a tag to an entity.
func (tm *TagManager) AddTag(e Entity, tag string) {
	tag = filterTag(tag)
	if tag == "" {
		return
	}

	// Initialize entity's tag set if needed
	if tm.entityTags[e] == nil {
		tm.entityTags[e] = make(map[string]bool)
	}

	// Skip if already has this tag (prevents duplicates)
	if tm.entityTags[e][tag] {
		return
	}

	// Add to entity's tags
	tm.entityTags[e][tag] = true

	// Update reverse index
	if tm.tagIndex[tag] == nil {
		tm.tagIndex[tag] = make(map[Entity]bool)
	}
	tm.tagIndex[tag][e] = true
}

// RemoveTag removes a tag from an entity.
func (tm *TagManager) RemoveTag(e Entity, tag string) {
	tag = filterTag(tag)
	if tag == "" {
		return
	}

	// Remove from entity's tags
	delete(tm.entityTags[e], tag)

	// Remove from reverse index
	if tm.tagIndex[tag] != nil {
		delete(tm.tagIndex[tag], e)
	}
}

// HasTag checks if an entity has a specific tag.
func (tm *TagManager) HasTag(e Entity, tag string) bool {
	tag = filterTag(tag)
	return tm.entityTags[e][tag]
}

// GetTags returns all tags for an entity.
func (tm *TagManager) GetTags(e Entity) []string {
	tags := tm.entityTags[e]
	if tags == nil {
		return nil
	}
	result := make([]string, 0, len(tags))
	for tag := range tags {
		result = append(result, tag)
	}
	return result
}

// GetEntitiesByTag returns all entities with the specified tag (O(1) lookup).
func (tm *TagManager) GetEntitiesByTag(tag string) []Entity {
	tag = filterTag(tag)
	entities := tm.tagIndex[tag]
	if entities == nil {
		return nil
	}
	result := make([]Entity, 0, len(entities))
	for e := range entities {
		result = append(result, e)
	}
	return result
}

// RemoveAllTags removes all tags from an entity. Called during Cleanup.
func (tm *TagManager) RemoveAllTags(e Entity) {
	tags := tm.entityTags[e]
	if tags == nil {
		return
	}
	// Remove from all reverse indexes
	for tag := range tags {
		if tm.tagIndex[tag] != nil {
			delete(tm.tagIndex[tag], e)
		}
	}
	// Remove entity's tag set
	delete(tm.entityTags, e)
}

// filterTag normalizes: lowercase, keeps only a-z, 0-9, underscore
func filterTag(s string) string {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]

		// Convert uppercase to lowercase
		if c >= 'A' && c <= 'Z' {
			c = c + 32
		}

		// Keep only valid characters
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' {
			result = append(result, c)
		}
	}
	return string(result)
}
