# ember2D

Experimental 2D Game Engine written in Go, using Ebiten for rendering and a browser-based visual logic editor inspired by Unreal Blueprints.

This project is early-stage, experimental, and under active development.

---

## Getting Started

ember2D consists of two separate executables:

1. **Runtime** — The Ebiten game engine runner  
2. **Editor** — A small web server that serves the visual editor UI

Both are located in the `cmd/` directory.

---

## ▶ Running the Runtime (Ebiten game window)

This launches a blank 640×480 Ebiten window.

```bash
go run ./cmd/ember2d-runtime

```
---

## ▶ Running the Editor (Visual Logic Editor)
This launches a web server at `http://localhost:8080` that serves the visual logic editor.

```bash
go run ./cmd/ember2d-editor
```

Once running, open your browser and navigate to:
```
http://localhost:8080
```

## Project Structure
The project is organized as follows:

<pre>
ember2D/
│
├── cmd/
│   ├── ember2d-editor/      # Editor server (web UI)
│   └── ember2d-runtime/     # Game runtime (Ebiten)
│
├── internal/                # Engine internals (private packages)
│   ├── engine/              # core, entities, components, behaviors
│   ├── editor/              # editor helpers (future)
│   └── util/                # utilities, JSON helpers
│
├── web/                     # Editor UI (HTML/CSS/JS)
│   ├── index.html
│   ├── app.js
│   └── style.css
│
├── config/                  # example rules/level configs
├── docs/                    # technical documentation
├── examples/                # future examples
│
├── go.mod
└── README.md

</pre>

## Technical details, engine architecture, and design decisions:
```
/docs/technical_analysis.md
```

## Future Work

Visual rule cards (event → condition → action)

Behavior compiler (JSON → Go structs)

Entity inspector in the editor

Real canvas level editor

Runtime hot-reload of JSON configs