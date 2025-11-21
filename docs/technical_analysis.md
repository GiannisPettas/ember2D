
# ember2D — Technical Analysis & Architecture Notes (Updated)

## Overview
ember2D is an experimental 2D game engine written in Go, designed around a hybrid visual-logic system inspired by Unreal Blueprints.  
It uses **Ebiten** as a low-level rendering/input backend, while implementing its own:

- entity system  
- behavior system  
- event-driven runtime  
- browser-based visual editor  
- JSON → runtime compiler  

The goal is to offer a lightweight but flexible framework for building 2D game logic visually.

---

# 🧩 Engine Architecture

ember2D follows a layered architecture:

```
entity   → basic data (no dependencies)
core     → event system, context, world reference
behavior → triggers, conditions, actions, dispatcher
runtime  → Ebiten game loop, rendering
editor   → browser-based UI for visual rule creation
```

These layers form a **dependency flow**:

```
entity → core → behavior → runtime
editor (separate system)
```

No cycles.

---

## 1. Core Layer (`internal/engine/core`)
The **core** package contains the engine’s fundamental low-level systems:

### ✔ Event System  
A minimal, dependency-free event model:

```go
type Event struct {
    Type    EventType
    A       string
    B       string
    Payload map[string]any
}
```

Events are “dumb data”.  
They do not depend on behavior, dispatcher, or runtime.

### ✔ Context  
A runtime execution context containing:

- the `World`
- the current `Event`

Passed into all actions and conditions to mutate state safely.

```go
type Context struct {
    World *entity.World
    Event Event
}
```

### ✔ Purpose of the Core Layer
- Provide minimal state and event abstraction  
- Stay **completely independent** from behavior logic  
- Avoid circular dependencies  
- Maintain stability and testability  

---

## 2. Entity System (`internal/engine/entity`)
A lightweight ECS-inspired model:

### ✔ Entity
A container with:
- ID
- Components (map[string]any)

### ✔ World
Manages all entities, and is shared by Context and Dispatcher.

World does **not** depend on behavior or runtime.

---

## 3. Behavior System (`internal/engine/behavior`)
The behavior system defines how game logic is expressed.

### ✔ Behavior
A behavior binds together:

- Trigger  
- Conditions  
- Actions  

This is the compiled form of visual rules coming from the editor.

### ✔ Trigger
Defines **when** the behavior should execute:

- "start"
- "collision"
- "timer"
- custom events
- optional entity filters (A/B matching)

### ✔ Conditions
Boolean checks:

```go
Evaluate(ctx *core.Context) bool
```

### ✔ Actions
State mutations:

```go
Execute(ctx *core.Context)
```

Examples:
- Damage entity  
- Teleport player  
- Spawn enemy  
- Change animation  
- Emit new event  

### ✔ Dependency rules
- Conditions + actions import **core**, NOT the dispatcher  
- Behaviors do NOT know about the dispatcher  
- Dispatcher does NOT import core cyclically  
- Behavior package has NO dependency on the runtime

---

## 4. Dispatcher (inside `behavior`)
The **Dispatcher** is the core of the event-driven runtime.

It:

- Receives events via `Emit`
- Stores them in a queue
- During `Update`, processes every queued event:
  - Matches against behavior triggers
  - Builds a `Context`
  - Evaluates conditions
  - Executes actions
  - Handles loop-back behavior

The dispatcher depends on:

- `entity.World`
- `core.Event`
- `behavior.Behavior`

but NOT the opposite direction.

---

## 5. Runtime (Ebiten Game Loop)
The runtime lives under `cmd/ember2d-runtime`.

Responsibilities:

- Initialize Ebiten window
- Create the World
- Create Dispatcher
- Load behaviors (from JSON)
- Load level entities (from JSON)
- Call dispatcher.Update() each frame
- Render entities (future)

Ebiten is used only for:
- window management
- rendering
- input

All game logic stays in the ember2D layers.

---

## 6. Loader (`internal/engine/loader` — planned)
Will load:

- level.json → entities  
- rules.json → triggers, conditions, actions  

And compile JSON into behavior objects.

---

# 🎛 Visual Logic System (Hybrid Design)

The editor does not use a node-graph like Unreal, but a “rule card” approach:

```
+------------------------+
| Event: OnCollision     |
| Entities: player enemy |
+------------------------+
| Conditions:            |
|   player.hp > 0        |
+------------------------+
| Actions:               |
|   Damage(player, 20)   |
|   Flash(player, red)   |
|   LoopBack             |
+------------------------+
```

Rules compile to Go structs.

### Loop models supported:
- Loop N times  
- Loop until condition  
- LoopBack (infinite unless conditions prevent)  

---

# 🌐 Editor Architecture

The browser-based visual editor:

- Served from `/web` via editor binary
- index.html → UI layout
- app.js → rule editing + exporting JSON
- style.css → dark theme

Will ultimately:

- Drag/drop rule cards  
- Configure triggers, conditions, actions  
- Save JSON configs  
- Live preview via WebSocket  

The editor is fully decoupled from the runtime.

---

# 🔮 Future Work

### Runtime
- Animation system  
- Collision system  
- Physics-lite (AABB only)  
- Tilemap support  
- Rendering pipeline  

### Editor
- Full drag-drop level editor  
- Inspector panel  
- Sprite/asset loader  
- JSON schema validation  

### Engine
- Action library expansion  
- Condition library expansion  
- Hot-reload of behaviors  

---

# 🧠 Development Notes

- Core MUST remain minimal to prevent cycles  
- Behavior should never import dispatcher or runtime  
- JSON must remain human-readable  
- Actions/Conditions should be small and composable  
- Event system will be minimal: nothing "smart", only data  
- Dispatcher owns all event logic  
- Runtime should remain thin: Ebiten integration only  
