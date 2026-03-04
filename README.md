# real-time-forum

A real-time forum built with Go and vanilla JavaScript, featuring authentication, posts, comments, and private messaging using WebSockets.

This project is a single-page application (SPA) powered by a Go backend and a JavaScript-driven frontend, with SQLite for persistent storage.

---

## 🚀 Features

### 🔐 Authentication

* User registration with:

  * Nickname
  * Age
  * Gender
  * First Name
  * Last Name
  * E-mail
  * Password (hashed with bcrypt)
* Login using:

  * Nickname + Password **or**
  * E-mail + Password
* Session management with cookies
* Logout available from any page
* Non-authenticated users can only access login/register view

---

### 📝 Posts & Comments

* Create posts with categories
* View posts in a feed display
* Click on a post to:

  * View full content
  * See related comments
* Add comments to posts
* Data stored in SQLite

---

### 💬 Private Messages (Real-Time)

* Real-time chat powered by WebSockets
* Online/offline user list:

  * Ordered by last message sent (like Discord)
  * Alphabetical order for new users with no messages
* Chat features:

  * Load last 10 messages by default
  * Infinite scroll to load 10 more (with throttle/debounce)
  * Real-time message delivery (no refresh required)
  * Message format includes:

    * Timestamp
    * Sender username
    * Message content
* Chat section always visible

---

## 🏗 Tech Stack

### Backend

* **Go**

  * HTTP server
  * Session handling
  * Goroutines & channels
* **WebSockets**

  * `gorilla/websocket`
* **Database**

  * `sqlite3`
* **Authentication**

  * `bcrypt`
* **UUID**

  * `gofrs/uuid` or `google/uuid`

### Frontend

* **HTML**

  * Single HTML file (SPA structure)
* **CSS**

  * Custom styling
* **Vanilla JavaScript**

  * DOM manipulation
  * Fetch API
  * WebSocket client
  * View switching (no page reloads)

---

## 📂 Architecture Overview

```
/backend
  ├── main.go
  ├── handlers/
  ├── websocket/
  ├── database/

/frontend
  ├── index.html
  ├── style.css
  ├── script.js
```

* Backend handles:

  * Routing
  * Authentication
  * Database access
  * WebSocket connections
* Frontend handles:

  * Dynamic rendering
  * Page transitions
  * WebSocket communication
  * Scroll optimization (throttle/debounce)

---

## 🧠 Learning Outcomes

This project reinforces:

* HTTP & sessions
* Cookies & authentication flows
* SQL and database design
* Go routines & channels
* WebSocket communication (Go + JS)
* DOM manipulation
* SPA architecture without frameworks
* Event optimization (throttling & debouncing)

---

## ⚙️ Setup

1. Install Go
2. Install SQLite
3. Clone the repository
4. Run:

```bash
go mod tidy
go run main.go
```

5. Open in browser:

```
http://localhost:8080
```

---

## 📌 Constraints

* Only one HTML file (SPA)
* No frontend frameworks (React, Vue, Angular, etc.)
* Allowed Go packages:

  * Standard library
  * gorilla/websocket
  * sqlite3
  * bcrypt
  * gofrs/uuid or google/uuid

y