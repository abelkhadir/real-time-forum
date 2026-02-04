package ws

import (
	"encoding/json"
	"log"
	"net/http"

	db "real/backend/database"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Message struct {
	To   string `json:"to"`
	From string `json:"from"`
	Msg  string `json:"msg"`
}

type Client struct {
	Username string
	Conn     *websocket.Conn
}

var clients = make(map[string]*Client)

func WebSocketsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthenticated"})
		return
	}

	username, err := db.GetUserBySession(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthenticated"})
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
	db.AddOnline(username)
	BroadcastUpdateContacts()

	defer func() {
		delete(clients, username)
		conn.Close()
		db.RemoveOnline(username)
		BroadcastUpdateContacts()
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

		db.SaveMessage(msg.From, msg.To, msg.Msg)

		if targetClient, ok := clients[msg.To]; ok {
			if err := targetClient.Conn.WriteJSON(msg); err != nil {
				log.Println("Write error:", err)
			}
		} else {
			log.Println("Target client not found:", msg.To)
		}
	}
}

func BroadcastOnlineUsers() {
	usernames := make([]string, 0, len(clients))
	for username := range clients {
		usernames = append(usernames, username)
	}

	msg := map[string]interface{}{
		"type":  "updatecontacts",
		"users": usernames,
	}

	for _, client := range clients {
		if err := client.Conn.WriteJSON(msg); err != nil {
			log.Println("Failed to send updatecontacts to", client.Username, err)
		}
	}
}
