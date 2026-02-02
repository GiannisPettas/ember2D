# Entity Manager Implementation Plan

**Goal:** Enhance the current entity system with querying, performance optimization (pooling), and clean entity creation patterns.

**Approach:** Step-by-step implementation where we discuss and code each feature together.

---

## Current State Analysis

### What We Have Now

```go
// internal/engine/entity/entity.go
type Entity struct {
    ID         string
    Components map[string]any
}

type World struct {
    Entities map[string]*Entity
}

func NewWorld() *World
func (w *World) AddEntity(e *Entity)
func (w *World) GetEntity(id string) *Entity
```

**Capabilities:**
- ✅ Store entities in a map
- ✅ Retrieve by ID
- ❌ No safe creation
- ❌ No deletion
- ❌ No querying
- ❌ No pooling

---

## Implementation Roadmap

We'll build this in **5 incremental steps**, each adding one major capability.

### Step 1: Safe Entity Creation with ID Generation

**Goal:** Add a factory method that creates entities with unique IDs automatically.

**What We'll Add:**
```go
type World struct {
    Entities map[string]*Entity
    nextID   int  // ← NEW: ID counter
}

// NEW: Factory method
func (w *World) CreateEntity(prefix string) *Entity {
    id := fmt.Sprintf("%s_%d", prefix, w.nextID)
    w.nextID++
    
    entity := &Entity{
        ID:         id,
        Components: make(map[string]any),
    }
    w.Entities[id] = entity
    return entity
}
```

**Usage Example:**
```go
// Before (manual, error-prone)
player := &Entity{ID: "player_1", Components: make(map[string]any)}
world.AddEntity(player)

// After (safe, automatic)
player := world.CreateEntity("player")  // Creates "player_0"
```

**Benefits:**
- ✅ No ID collisions
- ✅ Human-readable IDs for debugging
- ✅ Less boilerplate

**Discussion Points:**
- Should we keep `AddEntity()` for manual cases, or remove it?
- Do we want optional custom IDs? (e.g., `CreateEntity("player", "custom_id")`)

---

### Step 2: Tags for Entity Classification

**Goal:** Add tags to entities so we can categorize them (player, enemy, bullet, etc.).

**What We'll Add:**
```go
type Entity struct {
    ID         string
    Tags       []string           // ← NEW: Tags for classification
    Components map[string]any
}

// NEW: Tag management methods
func (e *Entity) AddTag(tag string)
func (e *Entity) HasTag(tag string) bool
func (e *Entity) RemoveTag(tag string)
```

**Usage Example:**
```go
player := world.CreateEntity("player")
player.AddTag("player")
player.AddTag("controllable")

enemy := world.CreateEntity("enemy")
enemy.AddTag("enemy")
enemy.AddTag("hostile")
```

**Benefits:**
- ✅ Categorize entities without complex inheritance
- ✅ Multiple tags per entity (flexible classification)
- ✅ Foundation for querying

**Discussion Points:**
- Should tags be added during creation? (builder pattern)
- Do we need tag validation/constants to prevent typos?

---

### Step 3: Querying System

**Goal:** Find entities by tags or component types.

**What We'll Add:**
```go
// NEW: Query methods
func (w *World) GetEntitiesByTag(tag string) []*Entity
func (w *World) GetEntitiesWithComponent(componentName string) []*Entity
func (w *World) GetAllEntities() []*Entity
```

**Implementation Strategy:**

**Option A: Simple Iteration** (Start here)
```go
func (w *World) GetEntitiesByTag(tag string) []*Entity {
    var result []*Entity
    for _, entity := range w.Entities {
        if entity.HasTag(tag) {
            result = append(result, entity)
        }
    }
    return result
}
```
- ✅ Simple, easy to understand
- ❌ O(N) - scans all entities every time

**Option B: Indexed Lookup** (Optimize later if needed)
```go
type World struct {
    Entities   map[string]*Entity
    tagIndex   map[string][]*Entity  // tag -> entities with that tag
    nextID     int
}
```
- ✅ O(1) lookup
- ❌ More complex, requires maintaining index

**Usage Example:**
```go
// Find all enemies
enemies := world.GetEntitiesByTag("enemy")
for _, enemy := range enemies {
    // Apply damage, AI logic, etc.
}

// Find all entities with health component
damageable := world.GetEntitiesWithComponent("health")
```

**Benefits:**
- ✅ Enables behavior patterns (damage all enemies, heal all allies)
- ✅ Supports systems that operate on entity groups
- ✅ Critical for AI (find nearest enemy, etc.)

**Discussion Points:**
- Start with simple iteration or build indexes from the start?
- Do we need complex queries? (e.g., "enemies with health > 50")

---

### Step 4: Safe Entity Deletion

**Goal:** Remove entities without causing crashes or race conditions.

**What We'll Add:**
```go
type World struct {
    Entities   map[string]*Entity
    toDelete   []string           // ← NEW: Deferred deletion queue
    nextID     int
}

// NEW: Deletion methods
func (w *World) DestroyEntity(id string)
func (w *World) Cleanup()  // Called at end of frame
```

**Implementation:**
```go
func (w *World) DestroyEntity(id string) {
    // Don't delete immediately - mark for deletion
    w.toDelete = append(w.toDelete, id)
}

func (w *World) Cleanup() {
    // Actually delete entities at safe time
    for _, id := range w.toDelete {
        delete(w.Entities, id)
    }
    w.toDelete = w.toDelete[:0]  // Clear queue
}
```

**Integration with Game Loop:**
```go
func (g *Game) Update() error {
    g.Dispatcher.Update()      // Process behaviors
    g.World.Cleanup()          // Delete dead entities
    return nil
}
```

**Benefits:**
- ✅ Safe: Deletion happens between frames
- ✅ No crashes from deleting during iteration
- ✅ Predictable timing

**Discussion Points:**
- Should we emit an event when entity is destroyed? (`ENTITY_DESTROYED`)
- Do we need "before destroy" callbacks for cleanup?

---

### Step 5: Entity Pooling for Performance

**Goal:** Reuse entity memory instead of allocating/deallocating constantly.

**What We'll Add:**
```go
type World struct {
    Entities     map[string]*Entity
    toDelete     []string
    entityPool   []*Entity         // ← NEW: Pre-allocated entities
    poolIndex    int               // ← NEW: Next available pool slot
    nextID       int
}

// NEW: Pool management
func (w *World) initPool(size int)
func (w *World) resetEntity(e *Entity)  // Clear entity for reuse
```

**Implementation:**
```go
func NewWorld() *World {
    w := &World{
        Entities:   make(map[string]*Entity),
        toDelete:   make([]string, 0),
        entityPool: make([]*Entity, 1000),  // Pre-allocate 1000 entities
    }
    
    // Pre-create entities
    for i := 0; i < 1000; i++ {
        w.entityPool[i] = &Entity{
            Components: make(map[string]any),
            Tags:       make([]string, 0, 4),
        }
    }
    return w
}

func (w *World) CreateEntity(prefix string) *Entity {
    var entity *Entity
    
    // Try to reuse from pool
    if w.poolIndex < len(w.entityPool) {
        entity = w.entityPool[w.poolIndex]
        w.poolIndex++
        w.resetEntity(entity)  // Clear previous data
    } else {
        // Pool exhausted - allocate new (rare)
        entity = &Entity{
            Components: make(map[string]any),
            Tags:       make([]string, 0, 4),
        }
    }
    
    entity.ID = fmt.Sprintf("%s_%d", prefix, w.nextID)
    w.nextID++
    w.Entities[entity.ID] = entity
    return entity
}

func (w *World) Cleanup() {
    for _, id := range w.toDelete {
        delete(w.Entities, id)
    }
    w.toDelete = w.toDelete[:0]
    w.poolIndex = 0  // Reset pool for next frame
}
```

**Benefits:**
- ✅ Zero allocations during gameplay (after warmup)
- ✅ No GC pauses
- ✅ Smooth 60 FPS even with thousands of bullets/particles

**Trade-offs:**
- ❌ Uses more RAM (pre-allocated pool)
- ❌ More complex code
- ✅ Worth it for high-frequency entities (bullets, particles)

**Discussion Points:**
- What pool size? (1000? 10000?)
- Should we have separate pools for different entity types?
- Do we need pool resizing if we run out?

---

## Bonus: Builder Pattern for Clean Creation

**Goal:** Make entity creation more readable and chainable.

**What We'll Add:**
```go
// NEW: Builder functions (functional options pattern)
type EntityOption func(*Entity)

func WithTag(tag string) EntityOption {
    return func(e *Entity) {
        e.AddTag(tag)
    }
}

func WithComponent(name string, value any) EntityOption {
    return func(e *Entity) {
        e.Components[name] = value
    }
}

// Enhanced CreateEntity
func (w *World) CreateEntity(prefix string, options ...EntityOption) *Entity {
    entity := w.createEntityInternal(prefix)
    
    // Apply all options
    for _, opt := range options {
        opt(entity)
    }
    
    return entity
}
```

**Usage Example:**
```go
// Before (verbose)
player := world.CreateEntity("player")
player.AddTag("player")
player.AddTag("controllable")
player.Components["position"] = Position{X: 100, Y: 50}
player.Components["health"] = 100

// After (clean, readable)
player := world.CreateEntity("player",
    WithTag("player"),
    WithTag("controllable"),
    WithComponent("position", Position{X: 100, Y: 50}),
    WithComponent("health", 100),
)
```

**Benefits:**
- ✅ Readable, declarative entity creation
- ✅ Easy to add new options without breaking existing code
- ✅ Common Go pattern (used in standard library)

---

## Testing Strategy

For each step, we'll:

1. **Implement** the feature in `entity.go`
2. **Test** it in `main.go` with a simple example
3. **Verify** the output/behavior
4. **Discuss** any issues or improvements
5. **Move to next step**

### Example Test for Step 1 (Creation):
```go
// In main.go
world := entity.NewWorld()

player := world.CreateEntity("player")
fmt.Printf("Created: %s\n", player.ID)  // Output: "Created: player_0"

enemy1 := world.CreateEntity("enemy")
enemy2 := world.CreateEntity("enemy")
fmt.Printf("Created: %s, %s\n", enemy1.ID, enemy2.ID)
// Output: "Created: enemy_1, enemy_2"
```

---

## Implementation Order

We'll proceed in this order:

1. ✅ **Step 1: Safe Creation** (Foundation)
2. ✅ **Step 2: Tags** (Classification)
3. ✅ **Step 3: Querying** (Finding entities)
4. ✅ **Step 4: Deletion** (Lifecycle completion)
5. ✅ **Step 5: Pooling** (Performance optimization)
6. ✅ **Bonus: Builder Pattern** (API polish)

Each step builds on the previous one, so we can stop at any point and still have a working system.

---

## Questions Before We Start

1. **Do you want to start with Step 1 (Safe Creation)?**
2. **Should we test each step in `main.go` or create a separate test file?**
3. **Any concerns or modifications to this plan?**

Let me know when you're ready to begin Step 1! 🚀
