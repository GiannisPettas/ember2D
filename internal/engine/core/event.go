package core

// EventType is a human-readable string representing the "kind" of event.
// Using strings (not iota ints) makes debugging, logging, and JSON import/export much easier.
type EventType string

// Event is the fundamental "message" structure in ember2D.
// Most game logic happens as a reaction to events.
// This is designed to be simple but extensible.
type Event struct {
	// Type defines the kind of event (e.g. "collision", "timer", "input").
	// Example: "collision", "timer", "pickup", "custom_event"
	Type EventType

	// A and B represent the main entities involved in the event.
	// Most game events (e.g., collisions, interactions) concern 1 or 2 entities.
	// If only one entity is relevant, B can be left blank ("").
	// Example: player collides with enemy (A="player", B="enemy").
	A string
	B string

	// Payload is a flexible map for extra data:
	// - Custom fields (e.g. "damage": 20, "direction": "left")
	// - Extra entities for complex events (e.g. area of effect, explosions)
	// - Any data you want to pass along with the event
	//
	// This keeps the core event structure clean while allowing maximum extensibility.
	Payload map[string]any
}

// ---
// Example usages:
//
// 1. Simple collision between player and enemy:
//    Event{
//        Type: "collision",
//        A:    "player",
//        B:    "enemy1",
//        Payload: map[string]any{"damage": 15},
//    }
//
// 2. Timer event with only one entity:
//    Event{
//        Type: "timer",
//        A:    "spawner1",
//        B:    "",
//        Payload: map[string]any{"interval": 2.0},
//    }
//
// 3. Area-of-effect event affecting multiple entities:
//    Event{
//        Type: "explosion",
//        A:    "bomb1",
//        B:    "",
//        Payload: map[string]any{
//            "affected": []string{"enemy1", "enemy2", "crate1"},
//            "radius":   50,
//        },
//    }
//
// 4. Custom event from the editor:
//    Event{
//        Type: "powerup_collected",
//        A:    "player",
//        B:    "",
//        Payload: map[string]any{"powerup_type": "speed_boost"},
//    }
//---------------------------------------------------------------------------------------------
// Example: Accessing an entity attribute from an event
//
// Suppose your World has a method: GetEntityByID(id string) *Entity
//
// func (w *World) GetEntityByID(id string) *Entity {
//     return w.Entities[id]
// }
//
// Suppose your Entity has a Components map (or fields).
//
// // Inside an Action or Condition:
// func SlowDownPlayer(ctx *core.Context) {
//     player := ctx.World.GetEntityByID(ctx.Event.A)
//     if player != nil {
//         // Assume the speed is stored as a float64 in Components
//         if speed, ok := player.Components["speed"].(float64); ok {
//             player.Components["speed"] = speed * 0.8
//         }
//     }
// }
//
// This ensures you always get the *current* value for the entity,
// not a stale copy from when the event was emitted.
