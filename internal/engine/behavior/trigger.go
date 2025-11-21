package behavior

import "github.com/GiannisPettas/ember2D/internal/engine/core"

// Trigger defines when a behavior should run.
type Trigger struct {
	Type     string   // e.g. "start", "collision", "timer"
	Entities []string // optional: entity ids or roles
	Interval int      // used for timers (future)
}

func (t Trigger) Matches(ev core.Event) bool {
	// Type check
	if string(ev.Type) != t.Type {
		return false
	}

	// No entity filter → always matches
	if len(t.Entities) == 0 {
		return true
	}

	// Collision / pair events: check against A/B
	for _, e := range t.Entities {
		if ev.A == e || ev.B == e {
			return true
		}
	}

	return false
}
