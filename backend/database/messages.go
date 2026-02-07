package db

type Message struct {
	To   string `json:"to"`
	From string `json:"from"`
	Msg  string `json:"msg"`
}

func SaveMessage(from, to, message string) error {
	_, err := db.Exec(`INSERT INTO messages (from_username, to_username, content) VALUES (?, ?, ?)`, from, to, message)

	return err
}

func ReadMessages(from, to string) ([]Message, error) {
	rows, err := db.Query(`SELECT from_username, to_username, content FROM messages 
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
		if err := rows.Scan(&msg.From, &msg.To, &msg.Msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func ReadMessagesPaged(from, to string, limit, offset int) ([]Message, error) {
	rows, err := db.Query(`SELECT from_username, to_username, content FROM messages 
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
		if err := rows.Scan(&msg.From, &msg.To, &msg.Msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
