package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

type editorMessage struct {
	Action string         `json:"action"`
	Data   map[string]any `json:"data,omitempty"`
}

type serverMessage struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Dev-friendly default for local editor use.
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	send := make(chan serverMessage, 32)
	done := make(chan struct{})

	go func() {
		defer close(done)
		for msg := range send {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("websocket write failed: %v", err)
				return
			}
		}
	}()

	send <- serverMessage{
		Type: "connected",
		Data: map[string]any{
			"server_time": time.Now().UTC().Format(time.RFC3339),
		},
	}

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var in editorMessage
		if err := json.Unmarshal(payload, &in); err != nil {
			send <- serverMessage{
				Type: "error",
				Data: map[string]any{"message": "invalid JSON payload"},
			}
			continue
		}

		switch in.Action {
		case "ping":
			send <- serverMessage{
				Type: "pong",
				Data: map[string]any{
					"server_time": time.Now().UTC().Format(time.RFC3339),
				},
			}
		default:
			send <- serverMessage{
				Type: "error",
				Data: map[string]any{
					"message": "unknown action: " + in.Action,
				},
			}
		}
	}

	close(send)
	<-done
}

func main() {
	mux := http.NewServeMux()

	// Serve the "web" folder as static files.
	webDir := filepath.Join("web")
	fs := http.FileServer(http.Dir(webDir))
	mux.Handle("/", fs)
	mux.HandleFunc("/ws", handleWS)

	log.Println("ember2D Editor running at http://localhost:9000")
	log.Println("WebSocket endpoint: ws://localhost:9000/ws")
	log.Fatal(http.ListenAndServe(":9000", mux))
}
