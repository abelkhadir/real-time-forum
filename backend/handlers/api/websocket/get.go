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
	var username string
	if err != nil {
		username = "guest"
	} else {
		username, err = db.GetUserBySession(cookie.Value)
		if username == "" || err != nil {
			username = "guest"
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{Username: username, Conn: conn}
	clients[username] = client

	log.Println("Client connected:", username)

	if username != "guest" {
		db.AddOnline(username)
	}
	BroadcastContacts(username)

	defer func() {
		delete(clients, username)
		conn.Close()
		db.RemoveOnline(username)
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
		if msg.To != "" && msg.To != "guest" && msg.From != "guest" {
			db.AddNotification(msg.To, msg.From, msg.Msg)
		}

		if targetClient, ok := clients[msg.To]; ok {
			if err := targetClient.Conn.WriteJSON(map[string]interface{}{
				"type": "UpdateMessages",
				"from": msg.From,
				"msg":  msg.Msg,
			}); err != nil {
				log.Println("Write error:", err)
			}
		} else {
			log.Println("Target client ot found:", msg.To)
		}
	}
}

func BroadcastContacts(username string) {
	contacts, _ := db.GetContacts()
	for _, c := range clients {
		c.Conn.WriteJSON(map[string]interface{}{
			"type":     "UpdateContacts",
			"contacts": contacts,
			"username": c.Username,
		})
	}
}

func BroadcastPost(post db.Post) {
	for _, c := range clients {
		c.Conn.WriteJSON(map[string]interface{}{
			"type": "UpdatePosts",
			"post": post,
		})
	}
}
