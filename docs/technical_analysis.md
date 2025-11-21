# ember2D — Technical Analysis & Architecture Notes

## Overview
ember2D is an experimental 2D game engine designed around a hybrid visual-logic
system, inspired by Unreal Blueprints, but implemented with Go + Ebiten.

The engine is split into two major subsystems:

1. **Runtime Engine**  
   - Written in Go  
   - Uses Ebiten for rendering + input  
   - Loads JSON behavior definitions  
   - Converts JSON → Go structs → executable gameplay logic  
   - Handles events, collisions, behaviors, and entity state  

2. **Browser-Based Visual Editor**  
   - Served by a Go web server  
   - Built with HTML/CSS/JavaScript  
   - Allows visual creation of behaviors (“rule cards”)  
   - Outputs rules.json + level.json  
   - Communicates with editor server via WebSocket/HTTP  

---

## Engine Architecture
The engine uses a clean modular structure under `/internal/engine`:

### 1. Core
Defines basic interfaces and the runtime update loop.

### 2. Entity System
Entities contain:
- ID
- Sprite reference
- Role (player, enemy, object, etc.)
- Components (hitbox, health, etc.)
- Behaviors (linked to compiled behavior objects)

### 3. Component System
Components are optional features attached to entities.
Examples:
- Hitbox
- Health
- Transform (in future)
- Inventory (future)

### 4. Behavior System
Behaviors are made of:
- Trigger (event type)
- Conditions (logic blocks)
- Actions (execution blocks)
- Optional looping logic

Behavior execution flow:
Event → Trigger Match → Conditions → Actions → LoopBack


### 5. Event Dispatcher
Central system that:
- Receives events (`Emit`)
- Matches events to behaviors
- Executes matching behaviors
- Supports collision, timer, and custom events

### 6. Loader
Loads:
- level.json
- rules.json
Compiles them into Go structs:
- Behavior
- Condition implementations
- Action implementations

---

## Visual Logic System (Hybrid Design)
The editor does not show a full node graph like Unreal.
Instead, it uses “rule cards”:

<pre>
+------------------------+  
| Event: OnCollision     |  
| Entities: player enemy |  
+------------------------+  
| Conditions:            |  
| [player.hp > 0]        |  
+------------------------+  
| Actions:               |  
| [Damage player 20]     |  
| [Flash red 300ms]      |  
| [LoopBack]             |  
+------------------------+  
</pre>


### Loop Model
Supported loop types:
- Loop N times
- Loop until condition
- LoopBack (repeat actions)

---

## Web Editor Architecture
Files under `/web/`:

- index.html — main editor UI  
- app.js — logic for dragging cards, editing properties, exporting JSON  
- style.css — dark theme + layout  

Editor server:
- Serves static HTML/JS/CSS
- Provides endpoints for saving/loading rules
- Future: WebSocket for live preview

---

## Future Expansion Ideas
- Animation editor (sprite sheets)
- Tile-based level editor
- Node-based behaviors (advanced mode)
- Real-time sync editor → runtime
- Shader playground (fragment shaders via Ebiten)

---

## Development Notes
- The engine runtime should remain small and clean.
- Visual logic must compile to predictable Go behaviors.
- JSON should remain human-readable for debugging.
- The system must support rapid expansion of actions/conditions.
