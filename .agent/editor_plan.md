# Editor / Runtime Architecture Plan

## Goal

Build a visual node editor (browser) that communicates with the Go runtime via WebSocket, allowing users to create entities, attach components, and define behaviors visually — all exported as JSON.

---

## Architecture

```
┌──────────────────────────┐        WebSocket         ┌──────────────────────┐
│   Editor (Browser)       │ ←────── JSON ──────────→ │   Runtime (Go)       │
│                          │                          │                      │
│  LiteGraph.js            │   editor → runtime:      │  WebSocket server    │
│  - Node canvas           │   create/update/delete   │  JSON decoder        │
│  - Widgets (edit values) │                          │  ECS engine          │
│  - Export/Import JSON    │   runtime → editor:      │  ebiten render       │
│  - Type-safe connections │   state sync, errors     │                      │
│                          │                          │                      │
│  Served by Go HTTP       │                          │  cmd/ember2d-editor  │
└──────────────────────────┘                          └──────────────────────┘
```

---

## Βήματα Υλοποίησης

### Phase 1: WebSocket Foundation
- [ ] Go WebSocket server (`cmd/ember2d-editor/main.go`)
- [ ] Serve static HTML/JS files
- [ ] JSON message protocol (action-based)
- [ ] Basic ping/pong connection test

### Phase 2: JSON Protocol
- [ ] Define message types:
  - `create_entity` — tags + components
  - `update_component` — entity ID + component data
  - `delete_entity` — entity ID
  - `state_sync` — runtime → editor (entity list)
- [ ] JSON ↔ ComponentManager bridge in Go
- [ ] Handle unknown/dynamic component data

### Phase 3: Basic Web Editor
- [ ] HTML page with LiteGraph.js
- [ ] Custom node types:
  - `EntityNode` — create entity with tags
  - `PositionNode` — X, Y widgets
  - `DisplayNode` — width, height, RGB widgets
  - `VelocityNode` — X, Y widgets
- [ ] Connect nodes → generate JSON command
- [ ] Send to Go runtime via WebSocket
- [ ] See entity appear in ebiten window!

### Phase 4: State Sync
- [ ] Runtime sends entity state to editor (positions, tags)
- [ ] Editor shows live entity list
- [ ] Select entity → highlight in runtime
- [ ] Update entity from editor → updates in runtime

### Phase 5: Advanced Nodes
- [ ] Behavior nodes (triggers, conditions, actions)
- [ ] Connection type enforcement (only matching types connect)
- [ ] Cycle detection
- [ ] Flow priority manager
- [ ] Save/Load project (JSON file)

---

## JSON Protocol Examples

### Editor → Runtime
```json
{"action": "create_entity", "data": {
    "tags": ["enemy", "hostile"],
    "components": {
        "position": {"x": 100, "y": 200},
        "velocity": {"x": 2, "y": -1},
        "display": {"width": 20, "height": 20, "r": 255, "g": 50, "b": 50}
    }
}}

{"action": "update_component", "data": {
    "entity": 3,
    "component": "position",
    "values": {"x": 150, "y": 300}
}}

{"action": "delete_entity", "data": {"entity": 3}}
```

### Runtime → Editor
```json
{"type": "state_sync", "data": {
    "entities": [
        {"id": 0, "tags": ["player"], "alive": true,
         "components": {"position": {"x": 300, "y": 220}}}
    ]
}}

{"type": "error", "data": {"message": "Entity 99 not found"}}
```

---

## Tech Stack

| Component | Technology |
|-----------|------------|
| Runtime | Go + ebiten |
| WebSocket | gorilla/websocket |
| Editor UI | LiteGraph.js (vanilla JS) |
| HTTP Server | Go net/http (serves static files) |
| Protocol | JSON over WebSocket |

---

## File Structure

```
cmd/ember2d-editor/
├── main.go              # Go HTTP + WebSocket server
├── handler.go           # WebSocket message handler
└── static/
    ├── index.html       # Editor page
    ├── editor.js        # LiteGraph setup + custom nodes
    ├── nodes/
    │   ├── entity.js    # EntityNode
    │   ├── position.js  # PositionNode
    │   ├── display.js   # DisplayNode
    │   └── velocity.js  # VelocityNode
    └── lib/
        └── litegraph.js # LiteGraph library

internal/engine/
├── protocol/
│   ├── messages.go      # JSON message structs
│   └── handler.go       # Process incoming commands
```

---

## Checklist - Αύριο ξεκινάμε

- [ ] Phase 1: WebSocket server + static file serving
- [ ] Phase 2: JSON protocol + component bridge
- [ ] Phase 3: Basic LiteGraph editor with 4 node types
- [ ] Phase 4: Real-time state sync
- [ ] Phase 5: Advanced nodes + behaviors
