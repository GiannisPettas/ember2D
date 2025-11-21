package core

import (
	"github.com/GiannisPettas/ember2D/internal/engine/entity"
)

// Context is passed to conditions and actions during behavior execution.
type Context struct {
	World *entity.World
	Event Event
}

func NewContext(world *entity.World, ev Event) *Context {
	return &Context{
		World: world,
		Event: ev,
	}
}

func (c *Context) GetEntity(id string) *entity.Entity {
	return c.World.GetEntity(id)
}
