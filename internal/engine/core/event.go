package core

type EventType string

type Event struct {
	Type    EventType
	A       string
	B       string
	Payload map[string]any
}
