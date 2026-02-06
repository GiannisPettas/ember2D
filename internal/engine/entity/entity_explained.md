# Entity Manager Explained

This document explains the Entity Manager and what we've implemented in ember2D.

---

## What We've Built

### Entity Structure

```go
type Entity struct {
    ID         string             // Unique identifier (e.g., "enemy_0")
    tags       []string           // Private - managed by World
    Components map[string]any     // Game data (position, health, etc.)
}

type World struct {
    Entities  map[string]*Entity                // All entities
    tagIndex  map[string]map[string]*Entity     // tag -> entityID -> *Entity
    idCounter int                               // For generating unique IDs
}
```

**Key design:** Tags are private (`tags` lowercase). Only World can modify them, ensuring the index stays consistent.

---

## Features Implemented

### 1. Safe Entity Creation

```go
enemy := world.CreateEntity("enemy")
// Creates entity with:
// - ID: "enemy_0" (auto-generated, unique)
// - tags: ["enemy"] (auto-added from prefix)
// - Indexed automatically for O(1) queries
```

### 2. Tag Index for O(1) Queries

We maintain a reverse lookup index:

```go
tagIndex = {
    "enemy":  {"enemy_0": *Entity, "enemy_1": *Entity},
    "player": {"player_0": *Entity},
    "bullet": {"bullet_0": *Entity, "bullet_1": *Entity, ...},
}
```

**Why nested maps?**
- Outer map: O(1) lookup by tag
- Inner map: O(1) add/remove entity from tag

### 3. World-Level Tag Methods

All tag operations go through World to keep the index in sync:

```go
world.AddTag(entity, "hostile")     // Add tag + update index
world.RemoveTag(entity, "hostile")  // Remove tag + update index
world.HasTag(entity, "hostile")     // Check for tag
world.GetTags(entity)               // Get all tags (copy)
world.GetEntitiesByTag("enemy")     // Find all enemies - O(1) lookup!
```

---

## Performance Analysis

### Query Performance

| Operation | Old (no index) | New (with index) |
|-----------|----------------|------------------|
| Find all enemies | Loop 10,000 entities | Direct map lookup |
| Complexity | O(N) | O(1) + O(k) |

Where N = total entities, k = entities matching the tag.

### The GC Trade-off

`GetEntitiesByTag` currently allocates a new slice each call:

```go
result := make([]*Entity, 0, len(entityMap))
```

**Problem:** If called frequently (10×/frame × 60 FPS = 600 allocations/sec), GC works hard.

**Solution (Step 5):** We'll add:
- Callback pattern for zero allocations
- Slice pooling for reusing memory

---

## Current API

### World Methods

```go
world := entity.NewWorld()                     // Create new world
entity := world.CreateEntity("prefix")         // Create with auto-ID and auto-tag

world.AddTag(entity, "tag")                    // Add tag to entity
world.RemoveTag(entity, "tag")                 // Remove tag from entity
world.HasTag(entity, "tag")                    // Check if entity has tag
world.GetTags(entity)                          // Get all tags (copy)

world.GetEntity("entity_0")                    // Get by exact ID
world.GetEntitiesByTag("enemy")                // Get all with tag - O(1)!
```

### Tag Normalization

All tags are normalized automatically:
- Converted to lowercase
- Only `a-z`, `0-9`, `_` allowed
- Examples: `"Enemy"` → `"enemy"`, `"Boss-Type"` → `"bosstype"`

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────┐
│                  World                          │
│  ┌─────────────┐    ┌──────────────────────┐   │
│  │  Entities   │    │      tagIndex        │   │
│  │ map[id]*E   │    │ map[tag]map[id]*E    │   │
│  └─────────────┘    └──────────────────────┘   │
│         │                    │                  │
│         └────── sync'd ──────┘                  │
│                                                 │
│  AddTag() / RemoveTag() update BOTH maps       │
└─────────────────────────────────────────────────┘
```

---

## Why World Controls Tags (Not Entity)

If Entity could modify its own tags:
```go
entity.AddTag("boss")                // Tags updated
// BUT tagIndex NOT updated!
world.GetEntitiesByTag("boss")       // Returns nothing! BUG!
```

By making `tags` private and routing through World:
- Index always stays in sync
- No way to accidentally break consistency
- Single source of truth

---

## What's Next

Still to implement:

1. ~~**Querying** - `GetEntitiesByTag("enemy")`~~ ✅ Done!
2. **Safe Deletion** - `DestroyEntity()` with deferred cleanup
3. **Pooling** - Reuse entity memory, reduce GC pressure
4. **Zero-Allocation Queries** - Callback pattern for hot paths
