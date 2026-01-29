package conditions

import "github.com/GiannisPettas/ember2D/internal/engine/core"

// AlwaysTrue is a utility condition that never blocks execution.
type AlwaysTrue struct{}

func (c *AlwaysTrue) Evaluate(ctx *core.Context) bool {
	return true
}
