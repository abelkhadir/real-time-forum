package db

type Message struct {
	To        string `json:"to"`
	From      string `json:"from"`
	Msg       string `json:"msg"`
	CreatedAt string `json:"created_at"`
}

// SaveMessage stores a private message and returns its creation time.
func SaveMessage(from, to, message string) (string, error) {
	res, err := db.Exec(`INSERT INTO messages (from_username, to_username, content) VALUES (?, ?, ?)`, from, to, message)
	if err != nil {
		return "", err
	}

	insertedID, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	var createdAt string
	err = db.QueryRow(`SELECT created_at FROM messages WHERE id = ?`, insertedID).Scan(&createdAt)
	if err != nil {
		return "", err
	}

	return createdAt, nil
}

// ReadMessages loads the full message history between two users.
func ReadMessages(from, to string) ([]Message, error) {
	rows, err := db.Query(`SELECT from_username, to_username, content, created_at FROM messages 
        WHERE (from_username = ? AND to_username = ?) 
           OR (from_username = ? AND to_username = ?)
        ORDER BY created_at ASC;
    `, from, to, to, from) // Pass both pairs
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.From, &msg.To, &msg.Msg, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// ReadMessagesPaged loads a page of messages between two users.
func ReadMessagesPaged(from, to string, limit, offset int) ([]Message, error) {
	rows, err := db.Query(`SELECT from_username, to_username, content, created_at FROM messages 
        WHERE (from_username = ? AND to_username = ?) 
           OR (from_username = ? AND to_username = ?)
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?;
    `, from, to, to, from, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.From, &msg.To, &msg.Msg, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
