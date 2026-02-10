# Entity System Explained

This document explains the Entity System architecture in ember2D.

---

## Architecture Overview

ember2D uses an **Entity Component System (ECS)** with Go generics:

```
┌──────────────────────────────────────────────────────┐
│                      World                           │
│  ┌──────────────┐    ┌─────────────────────────┐    │
│  │    alive      │    │      TagManager          │    │
│  │ map[uint64]   │    │ entityTags + tagIndex    │    │
│  │   bool        │    │ (dual map, O(1) lookup)  │    │
│  └──────────────┘    └─────────────────────────┘    │
│                                                      │
│  CreateEntity() / DestroyEntity() / Cleanup()        │
└──────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────┐
│           ComponentManager[T] (external)             │
│  ┌──────────────────────────────────────────────┐   │
│  │  positions := ComponentManager[Position]      │   │
│  │  healths   := ComponentManager[Health]        │   │
│  │  velocities := ComponentManager[Velocity]     │   │
│  │  ...one manager per component type            │   │
│  └──────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────┘
```

---

## Entity

An Entity is just a number (`uint64`). It has no data — it's an ID.

```go
type Entity uint64

// Entity 0, 1, 2, 3... auto-incremented
player := world.CreateEntity("player")   // returns Entity(0)
enemy  := world.CreateEntity("enemy")    // returns Entity(1)
```

**Why uint64 instead of a struct?**
- Fast comparisons (single CPU instruction)
- Fast map lookups (integer hash vs string hash)
- Lightweight (8 bytes vs struct with pointers)

---

## World

World manages entity lifecycle — creation, destruction, and cleanup.

```go
world := entity.NewWorld()

// Create entities (with optional tags)
player := world.CreateEntity("player")
enemy  := world.CreateEntity("enemy", "hostile", "ai")
bullet := world.CreateEntity("bullet")

// Check if alive
world.IsAlive(player)   // true

// Mark for deletion (deferred — removed at end of frame)
world.DestroyEntity(enemy)
world.IsAlive(enemy)    // false (marked, not yet removed)

// Actually remove at end of frame
world.Cleanup()         // Now enemy is fully gone

// Count alive entities
world.EntityCount()     // 2 (player + bullet)
```

### Deferred Deletion

`DestroyEntity` does NOT remove immediately. It marks the entity, then `Cleanup()` removes it at end of frame. This prevents bugs from modifying collections during iteration.

```
Frame:  [Systems run] → [DestroyEntity marks] → [Cleanup removes]
```

---

## TagManager

Tags classify entities. Accessed via `world.Tags()`.

```go
// Add/Remove/Check tags
world.Tags().AddTag(player, "controllable")
world.Tags().HasTag(player, "controllable")  // true
world.Tags().RemoveTag(player, "controllable")

// Get all tags for an entity
world.Tags().GetTags(player)  // ["player"]

// Find all entities with a tag — O(1) lookup!
enemies := world.Tags().GetEntitiesByTag("enemy")
```

### Dual Index (O(1) both directions)

```
entityTags:  entity → set of tags       (What tags does entity 5 have?)
tagIndex:    tag    → set of entities    (Which entities are "enemy"?)
```

Both maps update together — always in sync.

### Tag Normalization

All tags are automatically normalized:
- Uppercase → lowercase: `"Player"` → `"player"`
- Special chars removed: `"Boss-Type"` → `"bosstype"`
- Underscore kept: `"ai_enabled"` → `"ai_enabled"`
- Empty/invalid ignored: `""`, `"!!!"` → skipped

---

## ComponentManager[T]

Generic, type-safe component storage. Lives in `internal/engine/components/`.

```go
// Define your component types (plain structs)
type Position struct { X, Y float64 }
type Health   struct { Current, Max int }
type Velocity struct { X, Y float64 }

// Create a manager for each type
positions  := components.NewComponentManager[Position]()
healths    := components.NewComponentManager[Health]()
velocities := components.NewComponentManager[Velocity]()

// Attach components to entities
positions.Add(player, Position{X: 100, Y: 50})
healths.Add(player, Health{Current: 100, Max: 100})

// Read components — returns *T (pointer), nil if missing
pos := positions.Get(player)   // *Position
pos.X += 10                    // Modify directly via pointer

// Check/Remove
positions.Has(player)     // true
positions.Remove(player)  // detach component

// Iterate all entities with this component
velocities.Each(func(e entity.Entity, vel *Velocity) {
    if pos := positions.Get(e); pos != nil {
        pos.X += vel.X
        pos.Y += vel.Y
    }
})
```

### Why generics instead of map[string]any?

| | Old (`map[string]any`) | New (`ComponentManager[T]`) |
|---|---|---|
| Type safety | Runtime (crashes) | Compile-time (errors caught early) |
| Access | `entity.Components["pos"].(*Position)` | `positions.Get(entity)` |
| Performance | String hash + unboxing | uint64 hash, direct pointer |
| Boilerplate | Type assertion every time | Zero — compiler knows the type |

---

## Complete Example: Game Loop

```go
// Setup
world := entity.NewWorld()

positions  := components.NewComponentManager[Position]()
velocities := components.NewComponentManager[Velocity]()
healths    := components.NewComponentManager[Health]()

// Create entities
player := world.CreateEntity("player")
positions.Add(player, Position{X: 400, Y: 300})
healths.Add(player, Health{Current: 100, Max: 100})

for i := 0; i < 10; i++ {
    enemy := world.CreateEntity("enemy")
    positions.Add(enemy, Position{X: float64(i * 80), Y: 0})
    velocities.Add(enemy, Velocity{X: 0, Y: 2})
}

// Game loop (each frame)
// 1. Movement system
velocities.Each(func(e entity.Entity, vel *Velocity) {
    if pos := positions.Get(e); pos != nil {
        pos.X += vel.X
        pos.Y += vel.Y
    }
})

// 2. Cleanup dead enemies
for _, e := range world.Tags().GetEntitiesByTag("enemy") {
    if pos := positions.Get(e); pos != nil && pos.Y > 600 {
        world.DestroyEntity(e)
    }
}

// 3. End of frame cleanup
world.Cleanup()
```

---

## File Structure

```
internal/engine/
├── components/
│   └── manager.go       # ComponentManager[T] — generic storage
├── entity/
│   ├── entity.go        # Entity (uint64), World lifecycle
│   ├── tags.go          # TagManager with dual index
│   └── entity_test.go   # 21 tests, all passing
└── core/
    └── context.go       # Context for behavior system
```
