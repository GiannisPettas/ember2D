package core

import (
	"github.com/GiannisPettas/ember2D/internal/engine/entity"
)

// Context is passed to conditions and actions during behavior execution.
// It provides access to the World and the current Event.
type Context struct {
	World *entity.World
	Event Event
}

// NewContext creates a new Context.
func NewContext(world *entity.World, ev Event) *Context {
	return &Context{
		World: world,
		Event: ev,
	}
}
