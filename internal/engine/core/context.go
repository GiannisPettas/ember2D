package core

import (
	"github.com/GiannisPettas/ember2D/internal/engine/entity"
)

// Context is passed to conditions and actions during behavior execution.
//
// The Context provides access to the entire game state (World) and the current event being processed.
// This way, every action or condition can inspect and modify any entity, or query information about
// the specific event that triggered the behavior.
//
// The World field is a pointer, so any changes made to entities/components
// inside actions/conditions immediately update the real game state.
// The Event field contains the event being handled (collision, timer, etc.),
// including the main involved entities (A, B) and extra payload data.
//
// Example usage in an action:
//
//	func DamageEntity(ctx *core.Context) {
//	    target := ctx.World.GetEntity(ctx.Event.B)
//	    if target != nil {
//	        if hp, ok := target.Components["hp"].(float64); ok {
//	            target.Components["hp"] = hp - 10
//	        }
//	    }
//	}
type Context struct {
	World *entity.World // Reference to all entities & state in the game world
	Event Event         // The current event being processed (collision, timer, etc)
}

// Context bundles together everything needed to evaluate a condition or run an action
// when a behavior is triggered by an event.
//
// Why use a Context (instead of just passing Event or World)?
// - Centralizes all relevant state for modular actions/conditions.
// - Provides direct access to the full, up-to-date game world and the triggering event.
// - Avoids global variables and keeps code testable and extensible.
// - Ensures every rule always has access to the latest entity/component data.
// - Allows future expansion (e.g., add user, logger, event history, etc.) without breaking function signatures.
//
// Typical flow:
//  1. Dispatcher observes an Event.
//  2. Dispatcher creates a Context: ctx := NewContext(world, event)
//  3. Each condition/action gets ctx, and can safely read or modify any entity using ctx.World and info from ctx.Event.
func NewContext(world *entity.World, ev Event) *Context {
	return &Context{
		World: world, // Reference to all entities & state in the game world
		Event: ev,    // The current event being processed (collision, timer, etc)
	}
}

// GetEntity provides easy access to an entity by its ID,
// delegating the lookup to the World object.
//
// This is a convenience helper for actions/conditions so you don't have to write
// ctx.World.GetEntity(id) everywhere — just ctx.GetEntity(id).
//
// Typical usage in an action or condition:
//
//   player := ctx.GetEntity(ctx.Event.A)
//   if player != nil {
//       // ...modify player attributes, components, etc.
//   }
//
// By always fetching entities from the World (not from the Event or a copy),
// you guarantee that you are working with the latest, up-to-date state.

func (c *Context) GetEntity(id string) *entity.Entity {
	return c.World.GetEntity(id)
}

// Example usage:
//
// Suppose you want to slow down the player when a collision occurs.
// In your action or condition, you can use the context to fetch the entity:
//
//   func SlowDown(ctx *core.Context) {
//       player := ctx.GetEntity(ctx.Event.A) // or ctx.Event.B, depending on who
//       if player != nil {
//           // Access a component, e.g. "speed"
//           if speed, ok := player.Components["speed"].(float64); ok {
//               player.Components["speed"] = speed * 0.8
//           }
//       }
//   }
//
// This pattern guarantees that you always work with the latest entity state,
// even if other behaviors have modified it after the event was emitted.
//
