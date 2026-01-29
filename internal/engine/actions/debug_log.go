package actions

import (
	"fmt"
	"log"

	"github.com/GiannisPettas/ember2D/internal/engine/core"
)

// DebugLog is an action that prints a message to the console.
// It is useful for verifying that a behavior has triggered.
type DebugLog struct {
	Message string
}

func (a *DebugLog) Execute(ctx *core.Context) {
	// We can access the event that triggered this via ctx.Event
	// For now, just print the static message.
	log.Printf("[ACTION] DebugLog: %s (Event: %s)", a.Message, ctx.Event.Type)

	if ctx.Event.Payload != nil {
		fmt.Printf("\tPayload: %v\n", ctx.Event.Payload)
	}
}
