package db

// Comment represents a comment structure
type Comment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// GetCommentsByPost retrieves all comments for a specific post
func GetCommentsByPost(postID int) ([]Comment, error) {
	var comments []Comment

	rows, err := db.Query(`
		SELECT id, post_id, user_id, username, content, created_at
		FROM comments c
		WHERE post_id = ?
		ORDER BY created_at DESC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content,
			&c.CreatedAt); err != nil {
			continue
		}

		comments = append(comments, c)
	}

	return comments, nil
}

// GetComment retrieves a single comment by ID
func GetComment(commentID int) (Comment, error) {
	var c Comment

	row := db.QueryRow(`
		SELECT id, post_id, user_id, username, content,
		       0,
		       0,
		       created_at
		FROM comments c
		WHERE id = ?
	`, commentID)

	err := row.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content,
		&c.CreatedAt)
	if err != nil {
		return Comment{}, err
	}

	return c, nil
}

// InsertComment adds a new comment and returns the comment ID
func InsertComment(postID int, userID int, username, content string) (int64, error) {
	result, err := db.Exec(
		`INSERT INTO comments (post_id, user_id, username, content)
		 VALUES (?, ?, ?, ?)`,
		postID, userID, username, content,
	)
	if err != nil {
		return 0, err
	}

	// Update post comment count
	_, err = db.Exec(
		`UPDATE posts SET comments_num = comments_num + 1 WHERE id = ?`,
		postID,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// DeleteComment removes a comment and updates the post count
func DeleteComment(commentID int) error {
	// Get the post ID first
	var postID int
	err := db.QueryRow(`SELECT post_id FROM comments WHERE id = ?`, commentID).Scan(&postID)
	if err != nil {
		return err
	}

	// Delete the comment
	_, err = db.Exec(`DELETE FROM comments WHERE id = ?`, commentID)
	if err != nil {
		return err
	}

	// Update post comment count
	_, err = db.Exec(
		`UPDATE posts SET comments_num = comments_num - 1 WHERE id = ?`,
		postID,
	)
	return err
}
