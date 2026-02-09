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
	if !world.HasTag(player, "player") {
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

	world.AddTag(entity, "enemy")

	if !world.HasTag(entity, "enemy") {
		t.Error("Entity should have 'enemy' tag")
	}
}

func TestTagNormalization(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")

	world.AddTag(entity, "Player")
	world.AddTag(entity, "ENEMY")
	world.AddTag(entity, "Boss_Fight")

	if !world.HasTag(entity, "player") {
		t.Error("'Player' should normalize to 'player'")
	}
	if !world.HasTag(entity, "enemy") {
		t.Error("'ENEMY' should normalize to 'enemy'")
	}
	if !world.HasTag(entity, "boss_fight") {
		t.Error("'Boss_Fight' should normalize to 'boss_fight'")
	}
}

func TestTagFiltering(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")

	world.AddTag(entity, "Enemy-Type") // Hyphen should be removed
	world.AddTag(entity, "Hostile!!!") // Exclamations should be removed
	world.AddTag(entity, "AI_Enabled") // Underscore should be kept

	if !world.HasTag(entity, "enemytype") {
		t.Error("'Enemy-Type' should become 'enemytype'")
	}
	if !world.HasTag(entity, "hostile") {
		t.Error("'Hostile!!!' should become 'hostile'")
	}
	if !world.HasTag(entity, "ai_enabled") {
		t.Error("'AI_Enabled' should become 'ai_enabled'")
	}
}

func TestNoDuplicateTags(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")

	world.AddTag(entity, "player")
	world.AddTag(entity, "player")
	world.AddTag(entity, "PLAYER")
	// expecting 2: "test" (auto) + "player"
	if len(entity.tags) != 2 {
		t.Errorf("Expected 1 tag, got %d (duplicates not prevented)", len(entity.tags))
	}
}

func TestEmptyTagIgnored(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test") // Has "test" tag

	world.AddTag(entity, "")
	world.AddTag(entity, "!!!") // Becomes empty after filtering
	// Now expecting 1: just "test"
	if len(entity.tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(entity.tags))
	}
}

func TestRemoveTag(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")

	world.AddTag(entity, "enemy")
	world.AddTag(entity, "hostile")
	world.RemoveTag(entity, "hostile")

	if world.HasTag(entity, "hostile") {
		t.Error("Tag 'hostile' should have been removed")
	}
	if !world.HasTag(entity, "enemy") {
		t.Error("Tag 'enemy' should still exist")
	}
}

func TestHasTagFalse(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")

	world.AddTag(entity, "player")

	if world.HasTag(entity, "enemy") {
		t.Error("Entity should not have 'enemy' tag")
	}
}

// ============================================
// Normalization Tests for All Tag Functions
// ============================================
func TestHasTagNormalization(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	world.AddTag(entity, "player") // Stored as "player"
	// All these should find the tag
	if !world.HasTag(entity, "player") {
		t.Error("HasTag should find 'player'")
	}
	if !world.HasTag(entity, "Player") {
		t.Error("HasTag should normalize 'Player' to 'player'")
	}
	if !world.HasTag(entity, "PLAYER") {
		t.Error("HasTag should normalize 'PLAYER' to 'player'")
	}
	if !world.HasTag(entity, "PlAyEr") {
		t.Error("HasTag should normalize 'PlAyEr' to 'player'")
	}
}
func TestHasTagFilteringInput(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	world.AddTag(entity, "enemytype")  // Stored as "enemytype"
	world.AddTag(entity, "enemy_type") // Stored as "enemy_type"
	// HasTag should filter input the same way
	if !world.HasTag(entity, "Enemy-Type") {
		t.Error("HasTag should filter 'Enemy-Type' to 'enemytype'")
	}
	if !world.HasTag(entity, "enemy-type!!!") {
		t.Error("HasTag should filter 'enemy-type!!!' to 'enemytype'")
	}
	if !world.HasTag(entity, "enemy_type!?&") {
		t.Error("HasTag should filter 'enemy_type!?&' to 'enemy_type'")
	}
}
func TestRemoveTagNormalization(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	world.AddTag(entity, "enemy")   // Stored as "enemy"
	world.AddTag(entity, "hostile") // Stored as "hostile"
	// Remove using different cases
	world.RemoveTag(entity, "ENEMY")
	if world.HasTag(entity, "enemy") {
		t.Error("RemoveTag should remove 'enemy' when called with 'ENEMY'")
	}
	if !world.HasTag(entity, "hostile") {
		t.Error("'hostile' should still exist")
	}
}
func TestRemoveTagFilteringInput(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	world.AddTag(entity, "bosstype") // Stored as "bosstype"
	// Remove using unfiltered input
	world.RemoveTag(entity, "Boss-Type!!!")
	if world.HasTag(entity, "bosstype") {
		t.Error("RemoveTag should filter input and remove 'bosstype'")
	}
}
func TestAddTagNormalization(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	// Add with different cases - should all become the same tag
	world.AddTag(entity, "tEst")
	world.AddTag(entity, "player")
	world.AddTag(entity, "Player")
	world.AddTag(entity, "PLAYER")
	world.AddTag(entity, "PlAyEr")
	if len(entity.tags) != 2 {
		t.Errorf("Expected 2 tags (test + player), got %d", len(entity.tags))
	}
}

// ============================================
// Edge Case Tests
// ============================================
func TestRemoveTagThatDoesNotExist(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	world.AddTag(entity, "player")
	world.AddTag(entity, "friendly")
	// Try to remove a tag that doesn't exist
	world.RemoveTag(entity, "enemy")
	// Original tags should still be there
	if len(entity.tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(entity.tags))
	}
	if !world.HasTag(entity, "player") {
		t.Error("'player' tag should still exist")
	}
	if !world.HasTag(entity, "friendly") {
		t.Error("'friendly' tag should still exist")
	}
	if !world.HasTag(entity, "test") {
		t.Error("'test' tag should still exist")
	}
}
func TestRemoveTagFromEmptyTags(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	world.RemoveTag(entity, "test")
	// Entity has no tags - this should not panic
	world.RemoveTag(entity, "test")
	world.RemoveTag(entity, "anything")
	if len(entity.tags) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(entity.tags))
	}
}
func TestHasTagOnEmptyTags(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity("test")
	// Entity has no tags
	if world.HasTag(entity, "anything") {
		t.Error("HasTag should return false for empty tags")
	}
}

// ============================================
// Step 3 Tests: Querying
// ============================================
func TestGetEntitiesByTag(t *testing.T) {
	world := NewWorld()
	// Create different entity types
	world.CreateEntity("player")
	world.CreateEntity("enemy")
	world.CreateEntity("enemy")
	world.CreateEntity("enemy")
	world.CreateEntity("bullet")
	// Query enemies
	enemies := world.GetEntitiesByTag("enemy")
	if len(enemies) != 3 {
		t.Errorf("Expected 3 enemies, got %d", len(enemies))
	}
	// Query player
	players := world.GetEntitiesByTag("player")
	if len(players) != 1 {
		t.Errorf("Expected 1 player, got %d", len(players))
	}
}
func TestGetEntitiesByTagEmpty(t *testing.T) {
	world := NewWorld()
	world.CreateEntity("player")
	// Query non-existent tag
	ghosts := world.GetEntitiesByTag("ghost")
	if len(ghosts) != 0 {
		t.Errorf("Expected 0 ghosts, got %d", len(ghosts))
	}
}
func TestGetEntitiesByTagNormalization(t *testing.T) {
	world := NewWorld()
	world.CreateEntity("enemy")
	world.CreateEntity("enemy")
	// Query with different cases
	enemies := world.GetEntitiesByTag("ENEMY")
	if len(enemies) != 2 {
		t.Errorf("Expected 2 enemies with 'ENEMY' query, got %d", len(enemies))
	}
}

// ============================================
// Step 4 Tests: Safe Deletion
// ============================================

func TestDestroyEntity(t *testing.T) {
	world := NewWorld()

	enemy := world.CreateEntity("enemy")

	if !world.IsAlive(enemy) {
		t.Error("New entity should be alive")
	}

	world.DestroyEntity(enemy.ID)

	if world.IsAlive(enemy) {
		t.Error("Entity should be marked as not alive after DestroyEntity")
	}

	// Entity still exists until Cleanup
	if world.GetEntity(enemy.ID) == nil {
		t.Error("Entity should still exist before Cleanup")
	}
}

func TestCleanup(t *testing.T) {
	world := NewWorld()

	enemy := world.CreateEntity("enemy")
	world.DestroyEntity(enemy.ID)
	world.Cleanup()

	// Now entity should be gone
	if world.GetEntity(enemy.ID) != nil {
		t.Error("Entity should be removed after Cleanup")
	}

	if len(world.Entities) != 0 {
		t.Errorf("Expected 0 entities, got %d", len(world.Entities))
	}
}

func TestCleanupRemovesFromTagIndex(t *testing.T) {
	world := NewWorld()

	world.CreateEntity("enemy")
	enemy2 := world.CreateEntity("enemy")
	world.CreateEntity("enemy")

	world.DestroyEntity(enemy2.ID)
	world.Cleanup()

	enemies := world.GetEntitiesByTag("enemy")
	if len(enemies) != 2 {
		t.Errorf("Expected 2 enemies after cleanup, got %d", len(enemies))
	}
}

func TestDoubleDestroy(t *testing.T) {
	world := NewWorld()

	enemy := world.CreateEntity("enemy")

	// Destroy twice - should not panic or add duplicate
	world.DestroyEntity(enemy.ID)
	world.DestroyEntity(enemy.ID)

	if len(world.entitiesToDelete) != 1 {
		t.Errorf("Expected 1 in delete list, got %d", len(world.entitiesToDelete))
	}

	world.Cleanup()

	if len(world.Entities) != 0 {
		t.Errorf("Expected 0 entities, got %d", len(world.Entities))
	}
}

func TestDestroyNonExistent(t *testing.T) {
	world := NewWorld()

	// Should not panic
	world.DestroyEntity("nonexistent_99")

	if len(world.entitiesToDelete) != 0 {
		t.Errorf("Expected 0 in delete list, got %d", len(world.entitiesToDelete))
	}
}

func TestIsAliveNil(t *testing.T) {
	world := NewWorld()

	if world.IsAlive(nil) {
		t.Error("IsAlive(nil) should return false")
	}
}
