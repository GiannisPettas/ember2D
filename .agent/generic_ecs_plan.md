# Generic ECS Architecture - Implementation Plan

## Γιατί αλλάζουμε;

### Πρόβλημα με το `map[string]any`

```go
// ΤΩΡΑ - Κάθε φορά γράφεις αυτό το "σεντόνι":
if posRaw, ok := entity.Components["position"]; ok {
    if pos, ok := posRaw.(*Position); ok {
        pos.X += 10  // Επιτέλους ο πραγματικός κώδικας!
    }
}
```

**Προβλήματα:**
1. ❌ Boilerplate code - πολύς κώδικας για απλές λειτουργίες
2. ❌ Runtime errors - αν γράψεις `"pos"` αντί `"position"`, crash!
3. ❌ Αργό - type assertions σε κάθε frame

---

## Η Νέα Αρχιτεκτονική

### Concept 1: Entity = Απλά ένας αριθμός

```go
// ΠΡΙΝ
type Entity struct {
    ID         string           // "player_0"
    Components map[string]any   // Αργό, unsafe
}

// ΜΕΤΑ
type Entity uint64  // Απλά ένας αριθμός: 0, 1, 2, 3...
```

**Γιατί;**
- Τα `uint64` συγκρίνονται πιο γρήγορα από strings
- Μπορούν να χρησιμοποιηθούν ως array index (ακόμα πιο γρήγορα από map)
- Δεν χρειάζεται memory allocation

---

### Concept 2: ComponentManager[T] - Type-Safe Storage

Αντί να αποθηκεύουμε τα components μέσα στο Entity, τα αποθηκεύουμε σε **ξεχωριστούς managers**, έναν για κάθε τύπο:

```
┌─────────────────────────────────────────────────────────┐
│                         World                           │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌─────────────────────┐  ┌─────────────────────┐      │
│  │ positions           │  │ velocities          │      │
│  │ ComponentManager    │  │ ComponentManager    │      │
│  │ [Position]          │  │ [Velocity]          │      │
│  ├─────────────────────┤  ├─────────────────────┤      │
│  │ Entity 0 → (10, 20) │  │ Entity 0 → (5, 0)   │      │
│  │ Entity 1 → (50, 30) │  │ Entity 2 → (0, -2)  │      │
│  │ Entity 2 → (0, 100) │  │                     │      │
│  └─────────────────────┘  └─────────────────────┘      │
│                                                         │
│  ┌─────────────────────┐                               │
│  │ healths             │  (Entity 1 δεν έχει velocity  │
│  │ ComponentManager    │   - και αυτό είναι OK!)       │
│  │ [Health]            │                               │
│  ├─────────────────────┤                               │
│  │ Entity 0 → 100      │                               │
│  │ Entity 1 → 50       │                               │
│  └─────────────────────┘                               │
└─────────────────────────────────────────────────────────┘
```

---

### Concept 3: Πώς λειτουργεί ο ComponentManager

```go
// Ο manager είναι generic - δουλεύει με οποιονδήποτε τύπο
type ComponentManager[T any] struct {
    data map[Entity]*T  // Entity ID -> Pointer σε Component
}

// Δημιουργία
func NewComponentManager[T any]() *ComponentManager[T] {
    return &ComponentManager[T]{
        data: make(map[Entity]*T),
    }
}

// Προσθήκη component σε entity
func (cm *ComponentManager[T]) Add(e Entity, component T) {
    cm.data[e] = &component
}

// Λήψη component (nil αν δεν υπάρχει)
func (cm *ComponentManager[T]) Get(e Entity) *T {
    return cm.data[e]
}

// Αφαίρεση
func (cm *ComponentManager[T]) Remove(e Entity) {
    delete(cm.data, e)
}

// Έλεγχος αν υπάρχει
func (cm *ComponentManager[T]) Has(e Entity) bool {
    _, exists := cm.data[e]
    return exists
}
```

**Το `[T any]` σημαίνει:** "Αυτός ο manager δουλεύει με οποιονδήποτε τύπο T"

---

### Concept 4: Χρήση στην πράξη

```go
// 1. Ορίζεις τα components (structs)
type Position struct { X, Y float64 }
type Velocity struct { X, Y float64 }
type Health struct { Current, Max int }

// 2. Δημιουργείς τους managers
positions := NewComponentManager[Position]()
velocities := NewComponentManager[Velocity]()
healths := NewComponentManager[Health]()

// 3. Δημιουργείς entities (απλοί αριθμοί)
player := Entity(0)
enemy := Entity(1)

// 4. Προσθέτεις components
positions.Add(player, Position{X: 100, Y: 50})
velocities.Add(player, Velocity{X: 5, Y: 0})
healths.Add(player, Health{Current: 100, Max: 100})

positions.Add(enemy, Position{X: 200, Y: 50})
healths.Add(enemy, Health{Current: 50, Max: 50})
// enemy δεν έχει velocity - είναι στατικός!

// 5. Game loop - ΚΑΘΑΡΟΣ ΚΩΔΙΚΑΣ!
pos := positions.Get(player)
vel := velocities.Get(player)
if pos != nil && vel != nil {
    pos.X += vel.X
    pos.Y += vel.Y
}
```

**Σύγκριση:**

| Παλιά | Νέα |
|-------|-----|
| `entity.Components["position"].(*Position)` | `positions.Get(entity)` |
| Runtime crash αν typo | Compile error αν typo |
| Type assertion κάθε φορά | Ο compiler ξέρει τον τύπο |

---

## Τι θα υλοποιήσουμε

### Files Structure (χρησιμοποιούμε υπάρχοντα folders)

```
internal/engine/
├── components/
│   └── manager.go      # ComponentManager[T] - ΝΕΟ
├── entity/
│   ├── entity.go       # Entity type (uint64) + World - REFACTOR
│   ├── tags.go         # TagManager - ΝΕΟ (extracted from entity.go)
│   └── entity_test.go  # Tests - UPDATE
└── ... (actions, behavior, conditions, core - μένουν ίδια)
```

### File 1: `entity.go`

```go
package entity

// Entity is just an ID - lightweight and fast
type Entity uint64

// World manages entities and their lifecycle
type World struct {
    nextID           Entity
    alive            map[Entity]bool
    entitiesToDelete []Entity
    tags             *TagManager  // Ξεχωριστός manager για tags
}

func NewWorld() *World { ... }
func (w *World) CreateEntity() Entity { ... }
func (w *World) DestroyEntity(e Entity) { ... }
func (w *World) IsAlive(e Entity) bool { ... }
func (w *World) Cleanup() { ... }
```

### File 2: `components/manager.go`

```go
package components

import "github.com/GiannisPettas/ember2D/internal/engine/entity"

// ComponentManager provides type-safe component storage
type ComponentManager[T any] struct {
    data map[entity.Entity]*T
}

func NewComponentManager[T any]() *ComponentManager[T] { ... }
func (cm *ComponentManager[T]) Add(e Entity, c T) { ... }
func (cm *ComponentManager[T]) Get(e Entity) *T { ... }
func (cm *ComponentManager[T]) Remove(e Entity) { ... }
func (cm *ComponentManager[T]) Has(e Entity) bool { ... }

// For iteration
func (cm *ComponentManager[T]) Each(fn func(Entity, *T)) { ... }
```

### File 3: `entity/tags.go`

```go
package entity

// TagManager handles entity tags with O(1) lookup
type TagManager struct {
    entityTags map[Entity]map[string]bool  // entity -> set of tags
    tagIndex   map[string]map[Entity]bool  // tag -> set of entities
}

func NewTagManager() *TagManager { ... }
func (tm *TagManager) AddTag(e Entity, tag string) { ... }
func (tm *TagManager) HasTag(e Entity, tag string) bool { ... }
func (tm *TagManager) RemoveTag(e Entity, tag string) { ... }
func (tm *TagManager) GetEntitiesByTag(tag string) []Entity { ... }
```

---

## Παράδειγμα: Πώς θα φαίνεται ένα Game

```go
package main

import "ember2d/internal/engine/entity"

// Define your components
type Position struct { X, Y float64 }
type Velocity struct { X, Y float64 }
type Sprite struct { TextureID string }

func main() {
    // Create world and component managers
    world := entity.NewWorld()
    positions := entity.NewComponentManager[Position]()
    velocities := entity.NewComponentManager[Velocity]()
    sprites := entity.NewComponentManager[Sprite]()
    
    // Create player
    player := world.CreateEntity()
    positions.Add(player, Position{X: 100, Y: 50})
    velocities.Add(player, Velocity{X: 0, Y: 0})
    sprites.Add(player, Sprite{TextureID: "player.png"})
    world.Tags().AddTag(player, "player")
    
    // Create enemies
    for i := 0; i < 10; i++ {
        enemy := world.CreateEntity()
        positions.Add(enemy, Position{X: float64(i * 50), Y: 200})
        velocities.Add(enemy, Velocity{X: -1, Y: 0})
        sprites.Add(enemy, Sprite{TextureID: "enemy.png"})
        world.Tags().AddTag(enemy, "enemy")
    }
    
    // Game loop
    for {
        // Movement system - clean!
        velocities.Each(func(e entity.Entity, vel *Velocity) {
            if pos := positions.Get(e); pos != nil {
                pos.X += vel.X
                pos.Y += vel.Y
            }
        })
        
        // Destroy dead enemies
        for _, e := range world.Tags().GetEntitiesByTag("enemy") {
            if pos := positions.Get(e); pos != nil && pos.X < 0 {
                world.DestroyEntity(e)
            }
        }
        
        world.Cleanup()
    }
}
```

---

## Checklist - Βήματα Υλοποίησης

- [x] Δημιουργία `components/manager.go` με το `ComponentManager[T]`
- [x] Refactor `entity/entity.go` - Entity = uint64
- [x] Δημιουργία `entity/tags.go` - ξεχωριστός TagManager
- [x] Update tests στο `entity/entity_test.go`
- [x] Verify με `go test` - 21/21 tests pass ✅
- [x] Update `core/context.go`
- [x] `go build ./...` - project compiles ✅

---

## Ερωτήσεις πριν προχωρήσουμε

1. **Κατάλαβες τη διαφορά** μεταξύ `map[string]any` και `ComponentManager[T]`?
2. **Entity = uint64**: OK? 
3. **Ξεχωριστοί managers**: Κάθε component type έχει τον δικό του manager. OK?
