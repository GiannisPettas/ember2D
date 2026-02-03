package main

import (
	"log"

	"github.com/GiannisPettas/ember2D/internal/engine/actions"
	"github.com/GiannisPettas/ember2D/internal/engine/behavior"
	"github.com/GiannisPettas/ember2D/internal/engine/conditions"
	"github.com/GiannisPettas/ember2D/internal/engine/core"
	"github.com/GiannisPettas/ember2D/internal/engine/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Dispatcher *behavior.Dispatcher
}

func (g *Game) Update() error {
	if g.Dispatcher != nil {
		g.Dispatcher.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Debug info could go here
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	// 1. Initialize World
	world := entity.NewWorld()

	// 2. Define a Test Behavior
	// Rule: When "START_GAME" event happens -> Log "Engine is running!"
	testBehavior := &behavior.Behavior{
		ID: "test_rule_1",
		Trigger: behavior.Trigger{
			Type: "START_GAME",
		},
		Conditions: []behavior.Condition{
			&conditions.AlwaysTrue{},
		},
		Actions: []behavior.Action{
			&actions.DebugLog{Message: "Engine is running!"},
		},
	}

	// 3. Initialize Dispatcher
	dispatcher := behavior.NewDispatcher(world, []*behavior.Behavior{testBehavior})

	// 4. Emit the initial event to kickstart logic
	dispatcher.Emit(core.Event{
		Type: "START_GAME",
		Payload: map[string]any{
			"param": "test_value",
		},
	})

	// 5. Run Game Loop
	ebiten.SetWindowTitle("ember2D Runtime")
	game := &Game{
		Dispatcher: dispatcher,
	}

	log.Println("Starting ember2D runtime...")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
