package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func WebSocketsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer func() {
		log.Println("Connection closed")
		conn.Close()
	}()

	for {
		// ReadMessage returns message type (text/binary) and raw bytes
		msgType, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Try to parse JSON first
		var msg map[string]any
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			// Not JSON, treat as plain text
			msg = map[string]any{"text": string(msgBytes)}
		}

		log.Println("Received:", msg)

		// Echo message back as JSON
		if err := conn.WriteJSON(msg); err != nil {
			log.Println("Write error:", err)
			break
		}

		// Optional: handle close frame gracefully
		if msgType == websocket.CloseMessage {
			log.Println("Client requested close")
			break
		}
	}
}
