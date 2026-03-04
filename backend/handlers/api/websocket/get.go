package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

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

var (
	clients   = make(map[string][]*Client)
	clientsMu sync.RWMutex
)

func WebSocketsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username, err := db.GetUserBySession(cookie.Value)
	if username == "" || err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{Username: username, Conn: conn}
	clientsMu.Lock()
	clients[username] = append(clients[username], client)
	clientsMu.Unlock()

	db.AddOnline(username)
	BroadcastContacts(username)

	defer func() {
		clientsMu.Lock()
		userClients := clients[username]
		for i, c := range userClients {
			if c == client {
				userClients = append(userClients[:i], userClients[i+1:]...)
				break
			}
		}
		if len(userClients) == 0 {
			delete(clients, username)
		} else {
			clients[username] = userClients
		}
		clientsMu.Unlock()
		conn.Close()
		if len(userClients) == 0 {
			db.RemoveOnline(username)
		}
		BroadcastContacts(username)
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
		if msg.To != "" {
			db.AddNotification(msg.To, msg.From, msg.Msg)
		}

		event := map[string]interface{}{
			"type": "UpdateMessages",
			"to":   msg.To,
			"from": msg.From,
			"msg":  msg.Msg,
		}

		clientsMu.RLock()
		targetClients := clients[msg.To]
		senderClients := clients[msg.From]
		clientsMu.RUnlock()

		if len(targetClients) > 0 {
			for _, targetClient := range targetClients {
				if err := targetClient.Conn.WriteJSON(event); err != nil {
					log.Println("Write error:", err)
				}
			}
		} else {
			log.Println("Target client not found:", msg.To)
		}

		if msg.From != msg.To {
			for _, senderClient := range senderClients {
				if err := senderClient.Conn.WriteJSON(event); err != nil {
					log.Println("Write error:", err)
				}
			}
		}
	}
}

func BroadcastContacts(username string) {
	contacts, _ := db.GetContacts()
	clientsMu.RLock()
	defer clientsMu.RUnlock()
	for _, userClients := range clients {
		for _, c := range userClients {
			c.Conn.WriteJSON(map[string]interface{}{
				"type":     "UpdateContacts",
				"contacts": contacts,
				"username": c.Username,
			})
		}
	}
}

func BroadcastPost(post db.Post) {
	clientsMu.RLock()
	defer clientsMu.RUnlock()
	for _, userClients := range clients {
		for _, c := range userClients {
			c.Conn.WriteJSON(map[string]interface{}{
				"type": "UpdatePosts",
				"post": post,
			})
		}
	}
}
