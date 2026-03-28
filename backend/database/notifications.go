package db

type Notification struct {
	From      string `json:"from"`
	Msg       string `json:"msg"`
	CreatedAt string `json:"created_at"`
}

// AddNotification stores an unread notification for a user.
func AddNotification(username, from, message string) error {
	_, err := db.Exec(
		`INSERT INTO notifications (username, from_username, content) VALUES (?, ?, ?)`,
		username,
		from,
		message,
	)
	return err
}

// ReadUnreadNotifications returns unread notifications for a user.
func ReadUnreadNotifications(username string, limit int) ([]Notification, error) {
	rows, err := db.Query(
		`SELECT from_username, content, created_at
		 FROM notifications
		 WHERE username = ? AND read_at IS NULL
		 ORDER BY created_at DESC
		 LIMIT ?`,
		username,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.From, &n.Msg, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

// MarkNotificationsRead marks all unread notifications for a user as read.
func MarkNotificationsRead(username string) error {
	_, err := db.Exec(
		`UPDATE notifications SET read_at = CURRENT_TIMESTAMP
		 WHERE username = ? AND read_at IS NULL`,
		username,
	)
	return err
}

// MarkNotificationsReadFrom marks unread notifications as read for one sender.
func MarkNotificationsReadFrom(username, from string) error {
	_, err := db.Exec(
		`UPDATE notifications SET read_at = CURRENT_TIMESTAMP
		 WHERE username = ? AND from_username = ? AND read_at IS NULL`,
		username,
		from,
	)
	return err
}
