package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"real/backend/handlers/api/auth/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Client struct {
	Username string
	Conn     *websocket.Conn
}

type Message struct {
	To   string `json:"to"`
	From string `json:"from"`
	Msg  string `json:"msg"`
}

var clients = make(map[string]*Client)

func WebSocketsHandler(w http.ResponseWriter, r *http.Request) {
	username, err := utils.GetUsernameFromSession(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{Username: username, Conn: conn}
	clients[username] = client
	log.Println("Client connected:", username)

	defer func() {
		delete(clients, username)
		conn.Close()
		log.Println("Client disconnected:", username)
	}()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Println("Invalid JSON:", err)
			continue
		}

		// Server assigns sender
		msg.From = username

		if targetClient, ok := clients[msg.To]; ok {
			if err := targetClient.Conn.WriteJSON(msg); err != nil {
				log.Println("Write error:", err)
			}
		} else {
			log.Println("Target client not found:", msg.To)
		}
	}
}
