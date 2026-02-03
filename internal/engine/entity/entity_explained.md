# Entity Manager Explained

This document explains the Entity Manager and what we've implemented in ember2D.

---

## What We've Built

### Entity Structure

```go
type Entity struct {
    ID         string             // Unique identifier (e.g., "enemy_0")
    Tags       []string           // Categories (e.g., ["enemy", "hostile"])
    Components map[string]any     // Game data (position, health, etc.)
}

type World struct {
    Entities  map[string]*Entity  // All entities
    idCounter int                 // For generating unique IDs
}
```

---

## Features Implemented

### 1. Safe Entity Creation

```go
enemy := world.CreateEntity("enemy")
// Creates entity with:
// - ID: "enemy_0" (auto-generated, unique)
// - Tags: ["enemy"] (auto-added from prefix)
// - Components: empty map (ready for data)
```

**Key behaviors:**
- `idCounter` ensures unique IDs across all entities
- Prefix is auto-added as tag for easy querying
- IDs are human-readable for debugging

### 2. Tag System

Tags categorize entities for querying and game logic.

```go
entity.AddTag("hostile")     // Add a tag
entity.AddTag("FLYING")      // Normalized to "flying"
entity.HasTag("hostile")     // Check for tag → true
entity.RemoveTag("hostile")  // Remove tag
```

**Tag Normalization:**
- All tags are converted to lowercase
- Only `a-z`, `0-9`, and `_` are allowed
- Invalid characters are removed
- Examples: `"Enemy"` → `"enemy"`, `"Boss-Type"` → `"bosstype"`

**Optimizations:**
- Single-pass `filterTag()` handles lowercase + filtering
- Early exit when entity has no tags

---

## Design Decisions

### Why Auto-Add Prefix as Tag?

When you create `world.CreateEntity("enemy")`, the engine automatically adds `"enemy"` as a tag.

**Benefits:**
- Entities are always queryable by their type
- Less boilerplate for game developers
- Hard to forget to tag entities

**Example:**
```go
// Without auto-tag (old way)
enemy := world.CreateEntity("enemy")
enemy.AddTag("enemy")  // Easy to forget!

// With auto-tag (current)
enemy := world.CreateEntity("enemy")
// Already has "enemy" tag!
```

### Why Normalize Tags?

Without normalization:
```go
entity.AddTag("Enemy")
entity.AddTag("enemy")  // Oops, different tag!
entity.HasTag("ENEMY")  // Returns false!
```

With normalization:
```go
entity.AddTag("Enemy")   // Stored as "enemy"
entity.AddTag("enemy")   // Duplicate, ignored
entity.HasTag("ENEMY")   // Returns true!
```

### ID vs Tags

| Concept | Purpose | Example |
|---------|---------|---------|
| **ID** | Unique identifier | `"player_0"`, `"enemy_42"` |
| **Tags** | Categories | `["enemy", "hostile", "boss"]` |

- **IDs** answer: "Which specific entity is this?"
- **Tags** answer: "What kind of entity is this?"

---

## Current API

### World Methods

```go
world := entity.NewWorld()                  // Create new world
entity := world.CreateEntity("prefix")      // Create entity with auto-ID and auto-tag
entity := world.GetEntity("entity_0")       // Get by exact ID
```

### Entity Methods

```go
entity.AddTag("tag")       // Add tag (normalized)
entity.HasTag("tag")       // Check for tag → bool
entity.RemoveTag("tag")    // Remove tag
```

---

## What's Next

Still to implement:

1. **Querying** - `GetEntitiesByTag("enemy")` to find groups
2. **Safe Deletion** - `DestroyEntity()` with deferred cleanup
3. **Pooling** - Reuse entity memory for performance
