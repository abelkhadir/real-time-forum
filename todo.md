### TODO — real-time-forum

#### Setup

* [x] Repo init (`/backend`, `/frontend`, `/db`)
* [ ] SQLite schema (users, sessions, posts, comments, messages)

#### Auth (Go + SQLite)

* [x] Register endpoint

  * [ ] Validate fields
  * [ ] `bcrypt` password hash
  * [ ] Unique nickname/email
* [ ] Login endpoint

  * [ ] Nickname **or** email + password
  * [ ] Session cookie (UUID)
* [ ] Logout endpoint

  * [ ] Invalidate session
* [ ] Auth middleware (HTTP + WS)

#### Backend (Go)

* [ ] HTTP server (routes)
* [ ] WebSocket hub

  * [ ] Register / unregister clients
  * [ ] Broadcast / direct messages
* [ ] Online/offline tracking
* [ ] Goroutines + channels for WS

#### Posts & Comments

* [ ] Create post (title, content, categories)
* [ ] List posts (feed)
* [ ] View single post
* [ ] Create comment
* [ ] Fetch comments on click

#### Private Messages (Core)

* [ ] Users list sidebar

  * [ ] Sort by last message
  * [ ] Fallback: alphabetical
* [ ] Open chat with user
* [ ] Send private message (WS)
* [ ] Receive message in real time (WS)
* [ ] Message format: `{date, from, content}`

#### Chat History

* [ ] Load last 10 messages
* [ ] Infinite scroll up

  * [ ] Throttle / debounce scroll
  * [ ] Load +10 messages
* [ ] Persist messages in SQLite

#### Frontend (JS – no frameworks)

* [ ] Single HTML file (SPA)
* [ ] Page switching via JS

  * [ ] login
  * [ ] register
  * [ ] feed
  * [ ] post view
  * [ ] chat
* [ ] WebSocket client
* [ ] DOM updates (messages, posts, users)

#### UI / UX

* [ ] Always-visible chat sidebar
* [ ] Online/offline indicator
* [ ] Unread message badge
* [ ] Basic responsive CSS

#### Security / Sanity

* [ ] Input validation
* [ ] SQL prepared statements
* [ ] WS auth check
* [ ] No message spoofing

#### Testing / Debug

* [ ] Multiple users WS test
* [ ] Refresh resilience (reconnect WS)
* [ ] Session expiration test

---

**If stuck / rabbit hole warning 🚨**

* WS bugs → log every WS event (`connect/send/recv/close`)
* Auth weirdness → dump cookies + session table
* Scroll spam → console log handler fire rate
* Ordering users → precompute `last_message_at` in SQL

If you want:

* SQL schema
* WS message protocol
* Minimal Go WS hub
* JS throttle/debounce snippet
