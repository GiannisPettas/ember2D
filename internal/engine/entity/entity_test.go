package entity

import "testing"

// ============================================
// Step 1 Tests: Entity Creation
// ============================================

func TestCreateEntity(t *testing.T) {
	world := NewWorld()

	player := world.CreateEntity("player")

	if player.ID != "player_0" {
		t.Errorf("Expected 'player_0', got '%s'", player.ID)
	}
	//Check that prefix is auto-added as tag
	if !player.HasTag("player") {
		t.Error("CreateEntity should auto-add prefix as tag")
	}
}

func TestCreateMultipleEntities(t *testing.T) {
	world := NewWorld()

	e1 := world.CreateEntity("enemy")
	e2 := world.CreateEntity("enemy")
	e3 := world.CreateEntity("bullet")

	if e1.ID != "enemy_0" {
		t.Errorf("Expected 'enemy_0', got '%s'", e1.ID)
	}
	if e2.ID != "enemy_1" {
		t.Errorf("Expected 'enemy_1', got '%s'", e2.ID)
	}
	if e3.ID != "bullet_2" {
		t.Errorf("Expected 'bullet_2', got '%s'", e3.ID)
	}
}

func TestEntitiesStoredInWorld(t *testing.T) {
	world := NewWorld()

	world.CreateEntity("player")
	world.CreateEntity("enemy")
	world.CreateEntity("enemy")

	if len(world.Entities) != 3 {
		t.Errorf("Expected 3 entities, got %d", len(world.Entities))
	}
}

func TestGetEntity(t *testing.T) {
	world := NewWorld()

	created := world.CreateEntity("player")
	retrieved := world.GetEntity("player_0")

	if retrieved == nil {
		t.Fatal("GetEntity returned nil")
	}

	if retrieved.ID != created.ID {
		t.Errorf("Retrieved entity ID '%s' doesn't match created '%s'", retrieved.ID, created.ID)
	}
}

func TestGetEntityNotFound(t *testing.T) {
	world := NewWorld()

	retrieved := world.GetEntity("nonexistent")

	if retrieved != nil {
		t.Error("Expected nil for nonexistent entity")
	}
}

// ============================================
// Step 2 Tests: Tags
// ============================================

func TestAddTag(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")

	entity.AddTag("enemy")

	if !entity.HasTag("enemy") {
		t.Error("Entity should have 'enemy' tag")
	}
}

func TestTagNormalization(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")

	entity.AddTag("Player")
	entity.AddTag("ENEMY")
	entity.AddTag("Boss_Fight")

	if !entity.HasTag("player") {
		t.Error("'Player' should normalize to 'player'")
	}
	if !entity.HasTag("enemy") {
		t.Error("'ENEMY' should normalize to 'enemy'")
	}
	if !entity.HasTag("boss_fight") {
		t.Error("'Boss_Fight' should normalize to 'boss_fight'")
	}
}

func TestTagFiltering(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")

	entity.AddTag("Enemy-Type") // Hyphen should be removed
	entity.AddTag("Hostile!!!") // Exclamations should be removed
	entity.AddTag("AI_Enabled") // Underscore should be kept

	if !entity.HasTag("enemytype") {
		t.Error("'Enemy-Type' should become 'enemytype'")
	}
	if !entity.HasTag("hostile") {
		t.Error("'Hostile!!!' should become 'hostile'")
	}
	if !entity.HasTag("ai_enabled") {
		t.Error("'AI_Enabled' should become 'ai_enabled'")
	}
}

func TestNoDuplicateTags(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")

	entity.AddTag("player")
	entity.AddTag("player")
	entity.AddTag("PLAYER")
	// expecting 2: "test" (auto) + "player"
	if len(entity.Tags) != 2 {
		t.Errorf("Expected 1 tag, got %d (duplicates not prevented)", len(entity.Tags))
	}
}

func TestEmptyTagIgnored(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test") // Has "test" tag

	entity.AddTag("")
	entity.AddTag("!!!") // Becomes empty after filtering
	// Now expecting 1: just "test"
	if len(entity.Tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(entity.Tags))
	}
}

func TestRemoveTag(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")

	entity.AddTag("enemy")
	entity.AddTag("hostile")
	entity.RemoveTag("hostile")

	if entity.HasTag("hostile") {
		t.Error("Tag 'hostile' should have been removed")
	}
	if !entity.HasTag("enemy") {
		t.Error("Tag 'enemy' should still exist")
	}
}

func TestHasTagFalse(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")

	entity.AddTag("player")

	if entity.HasTag("enemy") {
		t.Error("Entity should not have 'enemy' tag")
	}
}

// ============================================
// Normalization Tests for All Tag Functions
// ============================================
func TestHasTagNormalization(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	entity.AddTag("player") // Stored as "player"
	// All these should find the tag
	if !entity.HasTag("player") {
		t.Error("HasTag should find 'player'")
	}
	if !entity.HasTag("Player") {
		t.Error("HasTag should normalize 'Player' to 'player'")
	}
	if !entity.HasTag("PLAYER") {
		t.Error("HasTag should normalize 'PLAYER' to 'player'")
	}
	if !entity.HasTag("PlAyEr") {
		t.Error("HasTag should normalize 'PlAyEr' to 'player'")
	}
}
func TestHasTagFilteringInput(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	entity.AddTag("enemytype")  // Stored as "enemytype"
	entity.AddTag("enemy_type") // Stored as "enemy_type"
	// HasTag should filter input the same way
	if !entity.HasTag("Enemy-Type") {
		t.Error("HasTag should filter 'Enemy-Type' to 'enemytype'")
	}
	if !entity.HasTag("enemy-type!!!") {
		t.Error("HasTag should filter 'enemy-type!!!' to 'enemytype'")
	}
	if !entity.HasTag("enemy_type!?&") {
		t.Error("HasTag should filter 'enemy_type!?&' to 'enemy_type'")
	}
}
func TestRemoveTagNormalization(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	entity.AddTag("enemy")   // Stored as "enemy"
	entity.AddTag("hostile") // Stored as "hostile"
	// Remove using different cases
	entity.RemoveTag("ENEMY")
	if entity.HasTag("enemy") {
		t.Error("RemoveTag should remove 'enemy' when called with 'ENEMY'")
	}
	if !entity.HasTag("hostile") {
		t.Error("'hostile' should still exist")
	}
}
func TestRemoveTagFilteringInput(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	entity.AddTag("bosstype") // Stored as "bosstype"
	// Remove using unfiltered input
	entity.RemoveTag("Boss-Type!!!")
	if entity.HasTag("bosstype") {
		t.Error("RemoveTag should filter input and remove 'bosstype'")
	}
}
func TestAddTagNormalization(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	// Add with different cases - should all become the same tag
	entity.AddTag("tEst")
	entity.AddTag("player")
	entity.AddTag("Player")
	entity.AddTag("PLAYER")
	entity.AddTag("PlAyEr")
	if len(entity.Tags) != 2 {
		t.Errorf("Expected 2 tags (test + player), got %d", len(entity.Tags))
	}
}

// ============================================
// Edge Case Tests
// ============================================
func TestRemoveTagThatDoesNotExist(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	entity.AddTag("player")
	entity.AddTag("friendly")
	// Try to remove a tag that doesn't exist
	entity.RemoveTag("enemy")
	// Original tags should still be there
	if len(entity.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(entity.Tags))
	}
	if !entity.HasTag("player") {
		t.Error("'player' tag should still exist")
	}
	if !entity.HasTag("friendly") {
		t.Error("'friendly' tag should still exist")
	}
	if !entity.HasTag("test") {
		t.Error("'test' tag should still exist")
	}
}
func TestRemoveTagFromEmptyTags(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	entity.RemoveTag("test")
	// Entity has no tags - this should not panic
	entity.RemoveTag("test")
	entity.RemoveTag("anything")
	if len(entity.Tags) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(entity.Tags))
	}
}
func TestHasTagOnEmptyTags(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	// Entity has no tags
	if entity.HasTag("anything") {
		t.Error("HasTag should return false for empty tags")
	}
}
