
# ember2D — Best Practices for Go Game Development

## Memory & GC
- **Pre-allocate arrays/slices/maps** for sprites, entities, particles at game or level load.
- **Re-use objects** wherever possible (particle pools, bullet pools, etc.).
- **Avoid creating new structs/slices/maps inside your main game loop** if possible.
- **Minimize small allocations per frame** — batch, pre-calculate, or reuse.
- **Use object pools** for short-lived objects (particles, projectiles, etc.).

## Assets & Loading
- **Load all assets (images, sounds, rules, levels) at startup** and keep them in memory.
- **Use pure Go libraries** for asset loading to avoid FFI/cgo overhead.

## Game Logic
- **Always access/modify entity state through the World/Context**, not via globals or stale pointers.
- **Split entity state into small, composable components** (e.g. position, health, animation).
- **Maintain a clear architecture for event dispatching, behaviors, and actions.**

## Rendering
- **Reuse image/options structs (e.g. ebiten.DrawImageOptions) per frame** whenever possible.
- **Batch draw calls** if rendering many similar sprites in a loop.

## Profiling & Testing
- **Profile your game with Go's built-in profiler (`go tool pprof`)** if you notice lag or spikes.
- **Write unit tests** for your core game logic (behaviors, conditions, actions).

## Project Structure
- **Organize code modularly** (`internal/engine`, `cmd/editor`, `cmd/runtime`, etc.).
- **Keep documentation and comments up-to-date** as the engine grows.
- **Plan early for dynamic module/asset/level loading and unloading.**

---

## Extra Best Practices

### Networking
- **Use goroutines and channels** to handle networking logic (multiplayer, async events) in Go idiomatic ways.
- **Keep all networking code isolated** from your main game logic for easier debugging.
- **Serialize game state/packets using efficient formats** (e.g. JSON for debugging, binary for performance).
- **Validate all network data before applying changes to game state.**

### AI (Artificial Intelligence)
- **Keep AI logic stateless or store state in dedicated AI components.**
- **Pre-calculate AI decisions or paths outside the main game loop** if possible (background goroutines).
- **Avoid excessive allocations in AI (pathfinding, tree nodes, etc.)** — use pools where needed.
- **Profile complex AI for GC spikes or frame delays.**

### Sounds
- **Pre-load audio assets and keep references alive** during gameplay (don't stream from disk unless necessary).
- **Prefer Go-native sound libraries** (e.g. oto, beep) to avoid FFI/cgo.
- **Limit simultaneous sound playbacks to avoid spikes in resource usage.**
- **Batch or schedule sound triggers to avoid triggering the same effect many times in one frame.**

---
