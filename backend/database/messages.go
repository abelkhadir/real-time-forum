package db

type Message struct {
	To   string `json:"to"`
	From string `json:"from"`
	Msg  string `json:"msg"`
}

func SaveMessage(from, to, message string) error {
	db.Exec(`INSERT INTO messages (sender, receiver, content) VALUES (?, ?, ?)`, from, to, message)

	return nil
}

func ReadMessages(from, to string) ([]Message, error) {
	rows, err := db.Query(`SELECT m.* FROM messages m
		WHERE (m.from_username = ? AND m.to_username = ?)
		OR (m.from_username = ? AND m.to_username = ?)
		ORDER BY m.created_at ASC;
	`, from, to, to, from)
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
