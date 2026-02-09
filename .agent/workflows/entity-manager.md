---
description: Entity Manager Implementation - Step-by-step ECS enhancement
---

# Entity Manager Implementation Workflow

**Status:** In Progress  
**Last Updated:** 2026-02-09

## Implementation Checklist

### Step 1: Safe Entity Creation ✅
- [x] Add `idCounter` to World
- [x] Implement `CreateEntity(prefix string)` factory method
- [x] Auto-generate IDs in format `{prefix}_{counter}`
- [x] Auto-add prefix as tag for easy querying
- [x] Tests written and passing

### Step 2: Tags ✅
- [x] Add `tags []string` field to Entity
- [x] Implement tag methods: `AddTag`, `HasTag`, `RemoveTag`, `GetTags`
- [x] Add tag normalization (lowercase, filter special chars)
- [x] Prevent duplicate tags
- [x] Add `tagIndex` to World for O(1) queries
- [x] Tests written and passing

### Step 3: Querying ✅
- [x] Implement `GetEntitiesByTag(tag string)` with index lookup
- [x] Return empty slice for non-existent tags
- [x] Normalize tag input in queries
- [x] Tests written and passing

### Step 4: Safe Deletion ✅
- [x] Add `isAlive` field to Entity
- [x] Add `entitiesToDelete` slice to World
- [x] Implement `DestroyEntity(id string)` - marks for deletion
- [x] Implement `IsAlive(entity *Entity)` - check alive status
- [x] Implement `Cleanup()` - removes marked entities
- [x] Cleanup removes from tag index too
- [x] Handle double-destroy and non-existent entity gracefully
- [x] Tests written and passing

### Step 5: Entity Pooling ⏳
- [ ] Add `entityPool []*Entity` to World
- [ ] Add `poolIndex int` to World
- [ ] Pre-allocate pool in `NewWorld()`
- [ ] Modify `CreateEntity` to reuse from pool
- [ ] Implement `resetEntity()` to clear entity for reuse
- [ ] Modify `Cleanup()` to reset pool index
- [ ] Tests for pooling behavior

### Step 6 (Bonus): Builder Pattern
- [ ] Create `EntityOption` functional options type
- [ ] Implement `WithTag(tag string) EntityOption`
- [ ] Implement `WithComponent(name string, value any) EntityOption`
- [ ] Update `CreateEntity` to accept options variadic
- [ ] Tests for builder pattern

---

## Reference Files

- **Implementation:** `internal/engine/entity/entity.go`
- **Tests:** `internal/engine/entity/entity_test.go`
- **Plan Document:** `internal/engine/entity/entity_manager_plan.md`

---

## How to Run Tests

```powershell
// turbo
cd c:\Users\petta\ember2D
go test ./internal/engine/entity/... -v
```

---

## Notes

- Keep this file updated as we make progress
- Check off items as they're completed
- Add any issues or decisions in the notes section
