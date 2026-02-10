package components

import (
	"testing"

	"github.com/GiannisPettas/ember2D/internal/engine/entity"
)

// Test component types
type Health struct{ Current, Max int }

// ============================================
// Creation Tests
// ============================================

func TestNewComponentManager(t *testing.T) {
	positions := NewComponentManager[Position]()

	if positions == nil {
		t.Fatal("NewComponentManager returned nil")
	}
	if positions.Count() != 0 {
		t.Errorf("Expected 0 components, got %d", positions.Count())
	}
}

// ============================================
// Add / Get Tests
// ============================================

func TestAddAndGet(t *testing.T) {
	positions := NewComponentManager[Position]()
	player := entity.Entity(0)

	positions.Add(player, Position{X: 100, Y: 50})

	pos := positions.Get(player)
	if pos == nil {
		t.Fatal("Get returned nil")
	}
	if pos.X != 100 || pos.Y != 50 {
		t.Errorf("Expected {100, 50}, got {%f, %f}", pos.X, pos.Y)
	}
}

func TestGetReturnsPointer(t *testing.T) {
	positions := NewComponentManager[Position]()
	player := entity.Entity(0)

	positions.Add(player, Position{X: 10, Y: 20})

	// Modify via pointer
	pos := positions.Get(player)
	pos.X = 999

	// Should reflect the change
	pos2 := positions.Get(player)
	if pos2.X != 999 {
		t.Errorf("Expected X=999 after pointer modification, got %f", pos2.X)
	}
}

func TestGetNonExistent(t *testing.T) {
	positions := NewComponentManager[Position]()

	pos := positions.Get(entity.Entity(42))
	if pos != nil {
		t.Error("Get should return nil for non-existent entity")
	}
}

func TestAddOverwrites(t *testing.T) {
	positions := NewComponentManager[Position]()
	player := entity.Entity(0)

	positions.Add(player, Position{X: 10, Y: 20})
	positions.Add(player, Position{X: 99, Y: 88})

	pos := positions.Get(player)
	if pos.X != 99 || pos.Y != 88 {
		t.Errorf("Expected {99, 88} after overwrite, got {%f, %f}", pos.X, pos.Y)
	}
	if positions.Count() != 1 {
		t.Errorf("Expected count 1 after overwrite, got %d", positions.Count())
	}
}

// ============================================
// Has / Remove Tests
// ============================================

func TestHas(t *testing.T) {
	positions := NewComponentManager[Position]()
	player := entity.Entity(0)

	if positions.Has(player) {
		t.Error("Has should return false before Add")
	}

	positions.Add(player, Position{X: 1, Y: 2})

	if !positions.Has(player) {
		t.Error("Has should return true after Add")
	}
}

func TestRemove(t *testing.T) {
	positions := NewComponentManager[Position]()
	player := entity.Entity(0)

	positions.Add(player, Position{X: 1, Y: 2})
	positions.Remove(player)

	if positions.Has(player) {
		t.Error("Has should return false after Remove")
	}
	if positions.Get(player) != nil {
		t.Error("Get should return nil after Remove")
	}
	if positions.Count() != 0 {
		t.Errorf("Expected count 0 after Remove, got %d", positions.Count())
	}
}

func TestRemoveNonExistent(t *testing.T) {
	positions := NewComponentManager[Position]()

	// Should not panic
	positions.Remove(entity.Entity(999))
}

// ============================================
// Each Tests
// ============================================

func TestEach(t *testing.T) {
	positions := NewComponentManager[Position]()

	positions.Add(entity.Entity(0), Position{X: 10, Y: 0})
	positions.Add(entity.Entity(1), Position{X: 20, Y: 0})
	positions.Add(entity.Entity(2), Position{X: 30, Y: 0})

	count := 0
	totalX := 0.0

	positions.Each(func(e entity.Entity, pos *Position) {
		count++
		totalX += pos.X
	})

	if count != 3 {
		t.Errorf("Each should visit 3 entities, visited %d", count)
	}
	if totalX != 60 {
		t.Errorf("Expected totalX=60, got %f", totalX)
	}
}

func TestEachCanModify(t *testing.T) {
	positions := NewComponentManager[Position]()

	positions.Add(entity.Entity(0), Position{X: 10, Y: 20})
	positions.Add(entity.Entity(1), Position{X: 30, Y: 40})

	// Move all entities right by 5
	positions.Each(func(e entity.Entity, pos *Position) {
		pos.X += 5
	})

	pos0 := positions.Get(entity.Entity(0))
	pos1 := positions.Get(entity.Entity(1))

	if pos0.X != 15 {
		t.Errorf("Expected X=15, got %f", pos0.X)
	}
	if pos1.X != 35 {
		t.Errorf("Expected X=35, got %f", pos1.X)
	}
}

func TestEachEmpty(t *testing.T) {
	positions := NewComponentManager[Position]()

	count := 0
	positions.Each(func(e entity.Entity, pos *Position) {
		count++
	})

	if count != 0 {
		t.Errorf("Each on empty manager should visit 0, visited %d", count)
	}
}

// ============================================
// Count Tests
// ============================================

func TestCount(t *testing.T) {
	healths := NewComponentManager[Health]()

	healths.Add(entity.Entity(0), Health{Current: 100, Max: 100})
	healths.Add(entity.Entity(1), Health{Current: 50, Max: 100})
	healths.Add(entity.Entity(2), Health{Current: 75, Max: 100})

	if healths.Count() != 3 {
		t.Errorf("Expected count 3, got %d", healths.Count())
	}

	healths.Remove(entity.Entity(1))

	if healths.Count() != 2 {
		t.Errorf("Expected count 2 after remove, got %d", healths.Count())
	}
}

// ============================================
// Multiple Manager Tests
// ============================================

func TestMultipleManagers(t *testing.T) {
	positions := NewComponentManager[Position]()
	healths := NewComponentManager[Health]()

	player := entity.Entity(0)

	positions.Add(player, Position{X: 100, Y: 200})
	healths.Add(player, Health{Current: 100, Max: 100})

	pos := positions.Get(player)
	hp := healths.Get(player)

	if pos.X != 100 {
		t.Errorf("Expected X=100, got %f", pos.X)
	}
	if hp.Current != 100 {
		t.Errorf("Expected Current=100, got %d", hp.Current)
	}

	// Remove position but keep health
	positions.Remove(player)

	if positions.Has(player) {
		t.Error("Should not have position after remove")
	}
	if !healths.Has(player) {
		t.Error("Should still have health")
	}
}
