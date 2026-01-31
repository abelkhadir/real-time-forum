package db

import (
	"database/sql"
)

// Comment represents a comment structure
type Comment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	Likes     int    `json:"likes_count"`
	Dislikes  int    `json:"dislikes_count"`
	CreatedAt string `json:"created_at"`
}

// GetCommentsByPost retrieves all comments for a specific post
func GetCommentsByPost(postID int) ([]Comment, error) {
	var comments []Comment

	rows, err := db.Query(`
		SELECT id, post_id, user_id, username, content, 
		       IFNULL(likes_count, 0), IFNULL(dislikes_count, 0), created_at
		FROM comments
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
			&c.Likes, &c.Dislikes, &c.CreatedAt); err != nil {
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
		       IFNULL(likes_count, 0), IFNULL(dislikes_count, 0), created_at
		FROM comments
		WHERE id = ?
	`, commentID)

	err := row.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content,
		&c.Likes, &c.Dislikes, &c.CreatedAt)
	if err != nil {
		return Comment{}, err
	}

	return c, nil
}

// InsertComment adds a new comment and returns the comment ID
func InsertComment(postID int, userID int, username, content string) (int64, error) {
	result, err := db.Exec(
		`INSERT INTO comments (post_id, user_id, username, content, likes_count, dislikes_count)
		 VALUES (?, ?, ?, ?, 0, 0)`,
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

// UpdateCommentLike updates or creates a like/dislike for a comment
func UpdateCommentLike(userID int, commentID int, isLike bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if like already exists
	var existingIsLike sql.NullBool
	err = tx.QueryRow(
		`SELECT is_like FROM comment_likes WHERE user_id = ? AND comment_id = ?`,
		userID, commentID,
	).Scan(&existingIsLike)

	if err == sql.ErrNoRows {
		// Insert new like
		_, err = tx.Exec(
			`INSERT INTO comment_likes (user_id, comment_id, is_like) VALUES (?, ?, ?)`,
			userID, commentID, isLike,
		)
		if err != nil {
			return err
		}

		// Update like/dislike count
		if isLike {
			_, err = tx.Exec(
				`UPDATE comments SET likes_count = likes_count + 1 WHERE id = ?`,
				commentID,
			)
		} else {
			_, err = tx.Exec(
				`UPDATE comments SET dislikes_count = dislikes_count + 1 WHERE id = ?`,
				commentID,
			)
		}
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if existingIsLike.Valid {
		// Update existing like
		wasLike := existingIsLike.Bool

		if wasLike == isLike {
			// Remove the like/dislike
			_, err = tx.Exec(
				`DELETE FROM comment_likes WHERE user_id = ? AND comment_id = ?`,
				userID, commentID,
			)
			if err != nil {
				return err
			}

			if isLike {
				_, err = tx.Exec(
					`UPDATE comments SET likes_count = likes_count - 1 WHERE id = ?`,
					commentID,
				)
			} else {
				_, err = tx.Exec(
					`UPDATE comments SET dislikes_count = dislikes_count - 1 WHERE id = ?`,
					commentID,
				)
			}
			if err != nil {
				return err
			}
		} else {
			// Change from like to dislike or vice versa
			_, err = tx.Exec(
				`UPDATE comment_likes SET is_like = ? WHERE user_id = ? AND comment_id = ?`,
				isLike, userID, commentID,
			)
			if err != nil {
				return err
			}

			// Adjust counts
			if wasLike {
				_, err = tx.Exec(
					`UPDATE comments SET likes_count = likes_count - 1, dislikes_count = dislikes_count + 1 WHERE id = ?`,
					commentID,
				)
			} else {
				_, err = tx.Exec(
					`UPDATE comments SET dislikes_count = dislikes_count - 1, likes_count = likes_count + 1 WHERE id = ?`,
					commentID,
				)
			}
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
