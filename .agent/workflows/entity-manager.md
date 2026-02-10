---
description: Entity Manager Implementation - Step-by-step ECS enhancement
---

# Entity Manager Implementation Workflow

**Status:** Generic ECS Refactor Complete  
**Last Updated:** 2026-02-10

## Architecture

Entity system uses **Go generics** for type-safe component storage:
- `Entity` = `uint64` (not a struct)
- `ComponentManager[T]` for type-safe components (no more `map[string]any`)
- `TagManager` with O(1) reverse index lookup
- Deferred deletion via `DestroyEntity` + `Cleanup`

## Key Files

| File | Purpose |
|------|---------|
| `internal/engine/entity/entity.go` | Entity type (uint64), World lifecycle |
| `internal/engine/entity/tags.go` | TagManager with dual index |
| `internal/engine/components/manager.go` | Generic ComponentManager[T] |
| `internal/engine/entity/entity_test.go` | 21 tests, all passing |
| `internal/engine/core/context.go` | Context for behaviors |

## Implementation Checklist

### Step 1: Safe Entity Creation ✅
- [x] Entity = uint64
- [x] `CreateEntity(tags ...string)` with auto-tagging
- [x] Tests passing

### Step 2: Tags ✅
- [x] TagManager with `map[Entity]map[string]bool`
- [x] Reverse index `map[string]map[Entity]bool` for O(1) queries
- [x] Tag normalization (lowercase, filter special chars)
- [x] No duplicate tags
- [x] Tests passing

### Step 3: Querying ✅
- [x] `GetEntitiesByTag` with O(1) index lookup
- [x] Normalization in queries
- [x] Tests passing

### Step 4: Safe Deletion ✅
- [x] `DestroyEntity` marks for deletion
- [x] `Cleanup` removes at end of frame + cleans tag index
- [x] Double-destroy and non-existent handled
- [x] Tests passing

### Step 5: ComponentManager[T] ✅
- [x] Generic type-safe storage
- [x] `Add`, `Get`, `Remove`, `Has`, `Each`, `Count`
- [x] Compiles and integrates with entity package

### Next Steps ⏳
- [ ] Add ComponentManager tests
- [ ] Builder pattern (optional)
- [ ] Entity pooling (optional, profile first)
- [ ] Update `entity_explained.md` docs

## How to Run Tests

```powershell
// turbo
go test ./internal/engine/entity/... -v
```

## How to Build

```powershell
// turbo
go build ./...
```
