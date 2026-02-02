# Entity Manager Explained

This document explains the Entity Manager concept and how it relates to our ember2D engine architecture.

---

## What is an Entity Manager?

An **Entity Manager** is the system responsible for the **lifecycle** and **organization** of entities. Think of it as the "HR Department" for your game objects.

### Current State: Basic Storage

Your `World` is currently just a **container** (a map). It can:

- ✅ Store entities
- ✅ Retrieve entities by ID
- ❌ **Cannot** create entities safely
- ❌ **Cannot** delete entities
- ❌ **Cannot** query entities (e.g., "find all enemies")
- ❌ **Cannot** generate unique IDs
- ❌ **Cannot** handle entity pooling (for performance)

### Entity Manager: Full Lifecycle Management

An Entity Manager adds:

1. **Creation** - Factory methods to spawn entities with unique IDs
2. **Deletion** - Safe removal (handling references, cleanup)
3. **Querying** - Find entities by tag, component type, or criteria
4. **Iteration** - Loop through all/filtered entities efficiently
5. **Pooling** - Reuse dead entities instead of allocating new ones (performance!)

---

## The Architecture Layers

```
┌─────────────────────────────────────────┐
│         Game Loop (main.go)             │
│  "I need to spawn an enemy at (100,50)" │
└─────────────────┬───────────────────────┘
                  │
                  ↓
┌─────────────────────────────────────────┐
│      Entity Manager (entity.go)         │ ← The "Factory"
│  • Generates unique IDs                 │
│  • Creates entities with components     │
│  • Tracks active/inactive entities      │
│  • Provides query methods               │
└─────────────────┬───────────────────────┘
                  │
                  ↓
┌─────────────────────────────────────────┐
│         World (entity.go)               │ ← The "Storage"
│  • map[string]*Entity                   │
│  • Just holds the data                  │
└─────────────────────────────────────────┘
```

---

## Practical Example: Current vs. Manager

### **Current Way** (Manual, Error-Prone)

```go
// In main.go or some behavior
player := &entity.Entity{
    ID: "player_1",  // ← What if this ID already exists?
    Components: make(map[string]any),
}
player.Components["position"] = Position{X: 100, Y: 50}
world.AddEntity(player)
```

**Problems:**
- No ID collision detection
- Verbose component setup
- No cleanup when entity dies
- Hard to find "all enemies" or "all entities with health"

### **With Entity Manager** (Safe, Convenient)

```go
// Clean API
playerID := world.CreateEntity(
    WithTag("player"),
    WithPosition(100, 50),
    WithHealth(100),
)

// Later...
enemies := world.GetEntitiesByTag("enemy")
for _, enemy := range enemies {
    // Do something with each enemy
}

// Safe deletion
world.DestroyEntity(playerID)  // Handles cleanup automatically
```

---

## Key Design Decisions

### 1. ID Generation Strategy

#### Option A: Sequential IDs

```go
nextID := 0
func CreateEntity() string {
    id := fmt.Sprintf("entity_%d", nextID)
    nextID++
    return id
}
```

- ✅ Simple, fast
- ❌ IDs are predictable (not always bad)

#### Option B: UUID

```go
import "github.com/google/uuid"
id := uuid.New().String()
```

- ✅ Globally unique (good for networking)
- ❌ Longer strings, slightly slower

#### Option C: Hybrid (Recommended for games)

```go
func CreateEntity(prefix string) string {
    return fmt.Sprintf("%s_%d", prefix, nextID)
}
// Creates: "player_0", "enemy_1", "bullet_2"
```

- ✅ Human-readable for debugging
- ✅ Fast
- ✅ Type-safe (you can tell what it is by ID)

### 2. Deletion Strategy

#### Immediate Deletion (Simple but dangerous)

```go
delete(world.Entities, id)
```

- ❌ What if another system is iterating over entities?
- ❌ What if a behavior references this entity?

#### Deferred Deletion (Safe, recommended)

```go
type World struct {
    Entities      map[string]*Entity
    toDelete      []string  // ← Mark for deletion
}

func (w *World) DestroyEntity(id string) {
    w.toDelete = append(w.toDelete, id)
}

func (w *World) Cleanup() {  // Called at end of frame
    for _, id := range w.toDelete {
        delete(w.Entities, id)
    }
    w.toDelete = w.toDelete[:0]  // Clear list
}
```

- ✅ Safe: Deletion happens between frames
- ✅ Predictable: All systems finish processing first

### 3. Querying Strategy

#### Option A: Tags (Unity-style)

```go
type Entity struct {
    ID         string
    Tags       []string  // ["enemy", "flying", "boss"]
    Components map[string]any
}

// Query
enemies := world.GetEntitiesByTag("enemy")
```

#### Option B: Component Filtering (ECS-style)

```go
// Find all entities with Position AND Velocity
moving := world.GetEntitiesWithComponents("position", "velocity")
```

#### Option C: Both (Most flexible)

```go
// Find all flying enemies
results := world.Query(
    HasTag("enemy"),
    HasTag("flying"),
    HasComponent("health"),
)
```

---

## Recommended Implementation

For ember2D, we should build:

1. **Safe entity creation** with auto-generated IDs
2. **Deferred deletion** to avoid mid-frame crashes
3. **Tag-based querying** for finding groups of entities
4. **Component helpers** to make setup easier

This gives you 80% of the power with 20% of the complexity.

---

## Memory Management Considerations

### The Problem: Allocation Churn

```go
// BAD: Creates garbage every frame
for i := 0; i < 100; i++ {
    bullet := &Entity{...}  // ← New allocation!
    world.AddEntity(bullet)
}
```

At 60 FPS, this creates 6,000 allocations per second. The Garbage Collector will eventually pause the game to clean up.

### The Solution: Object Pooling

```go
type World struct {
    Entities   map[string]*Entity
    entityPool []*Entity  // Pre-allocated entities
    nextPoolIndex int
}

func (w *World) CreateEntity() *Entity {
    if w.nextPoolIndex < len(w.entityPool) {
        // Reuse existing entity
        e := w.entityPool[w.nextPoolIndex]
        w.nextPoolIndex++
        return e
    }
    // Pool exhausted, allocate new
    return &Entity{...}
}

func (w *World) Cleanup() {
    // Reset pool for next frame
    w.nextPoolIndex = 0
}
```

**Trade-off:**
- Uses more RAM (pre-allocated pool)
- Eliminates GC pauses
- Critical for bullets, particles, effects

---

## Next Steps

Once we implement the Entity Manager, we can:

1. Create entities through events (`SPAWN_ENEMY`, `CREATE_BULLET`)
2. Query entities in behaviors (e.g., "find nearest enemy")
3. Safely destroy entities when they die
4. Build systems that operate on entity groups (rendering, physics, AI)

This transforms the World from a simple map into a powerful game object database.
