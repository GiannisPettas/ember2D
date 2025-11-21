package behavior

import "github.com/GiannisPettas/ember2D/internal/engine/core"

// Behavior is a compiled rule: trigger + conditions + actions.
type Behavior struct {
	ID         string
	Trigger    Trigger
	Conditions []Condition
	Actions    []Action
}

// Condition is a logic block that can block a behavior from running.
type Condition interface {
	Evaluate(ctx *core.Context) bool
}

// Action is a logic block that mutates the world / entities.
type Action interface {
	Execute(ctx *core.Context)
}
