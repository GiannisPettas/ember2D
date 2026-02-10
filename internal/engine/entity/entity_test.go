package entity

import "testing"

// ============================================
// Entity Creation Tests
// ============================================

func TestCreateEntity(t *testing.T) {
	world := NewWorld()

	player := world.CreateEntity("player")

	if player != 0 {
		t.Errorf("Expected entity 0, got %d", player)
	}
	if !world.IsAlive(player) {
		t.Error("New entity should be alive")
	}
}

func TestCreateMultipleEntities(t *testing.T) {
	world := NewWorld()

	e0 := world.CreateEntity("enemy")
	e1 := world.CreateEntity("enemy")
	e2 := world.CreateEntity("bullet")

	if e0 != 0 || e1 != 1 || e2 != 2 {
		t.Errorf("Expected 0, 1, 2 - got %d, %d, %d", e0, e1, e2)
	}
}

func TestCreateEntityAutoTag(t *testing.T) {
	world := NewWorld()

	player := world.CreateEntity("player")

	if !world.Tags().HasTag(player, "player") {
		t.Error("CreateEntity should auto-add tag")
	}
}

func TestCreateEntityMultipleTags(t *testing.T) {
	world := NewWorld()

	enemy := world.CreateEntity("enemy", "hostile", "ai")

	if !world.Tags().HasTag(enemy, "enemy") {
		t.Error("Should have 'enemy' tag")
	}
	if !world.Tags().HasTag(enemy, "hostile") {
		t.Error("Should have 'hostile' tag")
	}
	if !world.Tags().HasTag(enemy, "ai") {
		t.Error("Should have 'ai' tag")
	}
}

func TestEntityCount(t *testing.T) {
	world := NewWorld()

	world.CreateEntity("player")
	world.CreateEntity("enemy")
	world.CreateEntity("enemy")

	if world.EntityCount() != 3 {
		t.Errorf("Expected 3 entities, got %d", world.EntityCount())
	}
}

// ============================================
// Tag Tests
// ============================================

func TestAddTag(t *testing.T) {
	world := NewWorld()
	e := world.CreateEntity()

	world.Tags().AddTag(e, "enemy")

	if !world.Tags().HasTag(e, "enemy") {
		t.Error("Entity should have 'enemy' tag")
	}
}

func TestTagNormalization(t *testing.T) {
	world := NewWorld()
	e := world.CreateEntity()

	world.Tags().AddTag(e, "Player")
	world.Tags().AddTag(e, "ENEMY")
	world.Tags().AddTag(e, "Boss_Fight")

	if !world.Tags().HasTag(e, "player") {
		t.Error("'Player' should normalize to 'player'")
	}
	if !world.Tags().HasTag(e, "enemy") {
		t.Error("'ENEMY' should normalize to 'enemy'")
	}
	if !world.Tags().HasTag(e, "boss_fight") {
		t.Error("'Boss_Fight' should normalize to 'boss_fight'")
	}
}

func TestTagFiltering(t *testing.T) {
	world := NewWorld()
	e := world.CreateEntity()

	world.Tags().AddTag(e, "Enemy-Type")
	world.Tags().AddTag(e, "Hostile!!!")
	world.Tags().AddTag(e, "AI_Enabled")

	if !world.Tags().HasTag(e, "enemytype") {
		t.Error("'Enemy-Type' should become 'enemytype'")
	}
	if !world.Tags().HasTag(e, "hostile") {
		t.Error("'Hostile!!!' should become 'hostile'")
	}
	if !world.Tags().HasTag(e, "ai_enabled") {
		t.Error("'AI_Enabled' should become 'ai_enabled'")
	}
}

func TestNoDuplicateTags(t *testing.T) {
	world := NewWorld()
	e := world.CreateEntity()

	world.Tags().AddTag(e, "player")
	world.Tags().AddTag(e, "player")
	world.Tags().AddTag(e, "PLAYER")

	tags := world.Tags().GetTags(e)
	if len(tags) != 1 {
		t.Errorf("Expected 1 tag, got %d (duplicates not prevented)", len(tags))
	}
}

func TestEmptyTagIgnored(t *testing.T) {
	world := NewWorld()
	e := world.CreateEntity()

	world.Tags().AddTag(e, "")
	world.Tags().AddTag(e, "!!!")

	tags := world.Tags().GetTags(e)
	if len(tags) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(tags))
	}
}

func TestRemoveTag(t *testing.T) {
	world := NewWorld()
	e := world.CreateEntity()

	world.Tags().AddTag(e, "enemy")
	world.Tags().AddTag(e, "hostile")
	world.Tags().RemoveTag(e, "hostile")

	if world.Tags().HasTag(e, "hostile") {
		t.Error("Tag 'hostile' should have been removed")
	}
	if !world.Tags().HasTag(e, "enemy") {
		t.Error("Tag 'enemy' should still exist")
	}
}

func TestHasTagNormalization(t *testing.T) {
	world := NewWorld()
	e := world.CreateEntity()
	world.Tags().AddTag(e, "player")

	if !world.Tags().HasTag(e, "Player") {
		t.Error("HasTag should normalize 'Player' to 'player'")
	}
	if !world.Tags().HasTag(e, "PLAYER") {
		t.Error("HasTag should normalize 'PLAYER' to 'player'")
	}
}

func TestRemoveTagNormalization(t *testing.T) {
	world := NewWorld()
	e := world.CreateEntity()
	world.Tags().AddTag(e, "enemy")

	world.Tags().RemoveTag(e, "ENEMY")

	if world.Tags().HasTag(e, "enemy") {
		t.Error("RemoveTag should normalize 'ENEMY' to 'enemy'")
	}
}

// ============================================
// Query Tests
// ============================================

func TestGetEntitiesByTag(t *testing.T) {
	world := NewWorld()

	world.CreateEntity("player")
	world.CreateEntity("enemy")
	world.CreateEntity("enemy")
	world.CreateEntity("enemy")
	world.CreateEntity("bullet")

	enemies := world.Tags().GetEntitiesByTag("enemy")
	if len(enemies) != 3 {
		t.Errorf("Expected 3 enemies, got %d", len(enemies))
	}

	players := world.Tags().GetEntitiesByTag("player")
	if len(players) != 1 {
		t.Errorf("Expected 1 player, got %d", len(players))
	}
}

func TestGetEntitiesByTagEmpty(t *testing.T) {
	world := NewWorld()
	world.CreateEntity("player")

	ghosts := world.Tags().GetEntitiesByTag("ghost")
	if len(ghosts) != 0 {
		t.Errorf("Expected 0 ghosts, got %d", len(ghosts))
	}
}

func TestGetEntitiesByTagNormalization(t *testing.T) {
	world := NewWorld()
	world.CreateEntity("enemy")
	world.CreateEntity("enemy")

	enemies := world.Tags().GetEntitiesByTag("ENEMY")
	if len(enemies) != 2 {
		t.Errorf("Expected 2 enemies with 'ENEMY' query, got %d", len(enemies))
	}
}

// ============================================
// Deletion Tests
// ============================================

func TestDestroyEntity(t *testing.T) {
	world := NewWorld()

	enemy := world.CreateEntity("enemy")

	if !world.IsAlive(enemy) {
		t.Error("New entity should be alive")
	}

	world.DestroyEntity(enemy)

	if world.IsAlive(enemy) {
		t.Error("Entity should be marked as not alive after DestroyEntity")
	}
}

func TestCleanup(t *testing.T) {
	world := NewWorld()

	enemy := world.CreateEntity("enemy")
	world.DestroyEntity(enemy)
	world.Cleanup()

	if world.IsAlive(enemy) {
		t.Error("Entity should not be alive after Cleanup")
	}

	if world.EntityCount() != 0 {
		t.Errorf("Expected 0 entities, got %d", world.EntityCount())
	}
}

func TestCleanupRemovesFromTagIndex(t *testing.T) {
	world := NewWorld()

	world.CreateEntity("enemy")
	enemy2 := world.CreateEntity("enemy")
	world.CreateEntity("enemy")

	world.DestroyEntity(enemy2)
	world.Cleanup()

	enemies := world.Tags().GetEntitiesByTag("enemy")
	if len(enemies) != 2 {
		t.Errorf("Expected 2 enemies after cleanup, got %d", len(enemies))
	}
}

func TestDoubleDestroy(t *testing.T) {
	world := NewWorld()

	enemy := world.CreateEntity("enemy")

	world.DestroyEntity(enemy)
	world.DestroyEntity(enemy)

	world.Cleanup()

	if world.EntityCount() != 0 {
		t.Errorf("Expected 0 entities, got %d", world.EntityCount())
	}
}

func TestDestroyNonExistent(t *testing.T) {
	world := NewWorld()

	// Should not panic
	world.DestroyEntity(Entity(999))
}
