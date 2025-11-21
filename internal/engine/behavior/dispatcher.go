package behavior

import (
	"github.com/GiannisPettas/ember2D/internal/engine/core"
	"github.com/GiannisPettas/ember2D/internal/engine/entity"
)

// Dispatcher receives events and routes them to matching behaviors.
type Dispatcher struct {
	World      *entity.World
	Behaviors  []*Behavior
	eventQueue []core.Event
}

func NewDispatcher(world *entity.World, behaviors []*Behavior) *Dispatcher {
	return &Dispatcher{
		World:     world,
		Behaviors: behaviors,
	}
}

// Emit adds an event to the queue.
func (d *Dispatcher) Emit(ev core.Event) {
	d.eventQueue = append(d.eventQueue, ev)
}

// Update processes all queued events.
func (d *Dispatcher) Update() {
	if len(d.eventQueue) == 0 {
		return
	}

	queue := d.eventQueue
	d.eventQueue = nil

	for _, ev := range queue {
		d.processEvent(ev)
	}
}

func (d *Dispatcher) processEvent(ev core.Event) {
	for _, b := range d.Behaviors {
		// 1. Trigger match
		if !b.Trigger.Matches(ev) {
			continue
		}

		// 2. Build context
		ctx := core.NewContext(d.World, ev)

		// 3. Conditions
		ok := true
		for _, cond := range b.Conditions {
			if !cond.Evaluate(ctx) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}

		// 4. Actions
		for _, act := range b.Actions {
			act.Execute(ctx)
		}

		// 5. LoopBack support (simple version)
		if b.Trigger.Type == "loop" {
			d.Emit(ev)
		}
	}
}
