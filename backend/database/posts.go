package db

import "strings"

func InsertPost(username, title, content string, categories []string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
        INSERT INTO posts(username, title, content)
        VALUES (?, ?, ?)
    `, username, title, content)
	if err != nil {
		return 0, err
	}

	postID, _ := res.LastInsertId()

	for _, cat := range categories {
		// Insert category if not exists
		_, err := tx.Exec(`INSERT OR IGNORE INTO categories(name) VALUES(?)`, cat)
		if err != nil {
			return 0, err
		}

		// Get category id
		var catID int64
		err = tx.QueryRow(`SELECT id FROM categories WHERE name = ?`, cat).Scan(&catID)
		if err != nil {
			return 0, err
		}

		// Link post -> category
		_, err = tx.Exec(`INSERT INTO post_categories(post_id, category_id) VALUES(?, ?)`, postID, catID)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return postID, nil
}

type Post struct {
	ID           int64
	Username     string
	Title        string
	Content      string
	Categories   []string
	Likes_num    int
	Dislikes_num int
	Comments_num int
	CreatedAt    string
}

func GetPosts(page, limit int) ([]Post, error) {
	offset := (page - 1) * limit

	rows, err := db.Query(`
    SELECT
      p.id,
      p.username,
      p.title,
      p.content,
      p.created_at,
      IFNULL(GROUP_CONCAT(c.name, ','), '') AS categories,
      (SELECT COUNT(*) FROM post_likes WHERE post_id = p.id AND is_like = 1) AS likes_num,
      (SELECT COUNT(*) FROM post_likes WHERE post_id = p.id AND is_like = 0) AS dislikes_num,
      (SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comments_num
    FROM posts p
    LEFT JOIN post_categories pc ON pc.post_id = p.id
    LEFT JOIN categories c ON c.id = pc.category_id
    GROUP BY p.id
    ORDER BY p.created_at DESC
    LIMIT ? OFFSET ?
    `, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var p Post
		var catStr string

		if err := rows.Scan(&p.ID, &p.Username, &p.Title, &p.Content, &p.CreatedAt, &catStr, &p.Likes_num, &p.Dislikes_num, &p.Comments_num); err != nil {
			return nil, err
		}

		if catStr != "" {
			p.Categories = strings.Split(catStr, ",")
		}

		posts = append(posts, p)
	}

	return posts, nil
}

func GetPost(id int) (Post, error) {
	var p Post
	var catStr string

	query := `
        SELECT
            p.id,
            p.username,
            p.title,
            p.content,
            p.created_at,
            IFNULL(GROUP_CONCAT(c.name, ','), '') AS categories,
            (SELECT COUNT(*) FROM post_likes WHERE post_id = p.id AND is_like = 1) AS likes_num,
            (SELECT COUNT(*) FROM post_likes WHERE post_id = p.id AND is_like = 0) AS dislikes_num,
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comments_num
        FROM posts p
        LEFT JOIN post_categories pc ON pc.post_id = p.id
        LEFT JOIN categories c ON c.id = pc.category_id
        WHERE p.id = ?
        GROUP BY p.id
    `

	row := db.QueryRow(query, id)

	err := row.Scan(
		&p.ID,
		&p.Username,
		&p.Title,
		&p.Content,
		&p.CreatedAt,
		&catStr,
		&p.Likes_num,
		&p.Dislikes_num,
		&p.Comments_num,
	)
	if err != nil {
		return Post{}, err
	}

	if catStr != "" {
		p.Categories = strings.Split(catStr, ",")
	}

	return p, nil
}

func UpdatePostLike(userID, postID int, isLike bool) error {

	var existing *bool

	err := db.QueryRow(`
        SELECT is_like FROM post_likes
        WHERE user_id = ? AND post_id = ?`,
		userID, postID).Scan(&existing)

	if err != nil {
		_, err = db.Exec(`
            INSERT INTO post_likes(user_id, post_id, is_like)
            VALUES (?, ?, ?)`,
			userID, postID, isLike)

		return err
	}

	if *existing == isLike {
		_, err = db.Exec(`
            DELETE FROM post_likes
            WHERE user_id = ? AND post_id = ?`,
			userID, postID)

		return err
	}

	_, err = db.Exec(`
        UPDATE post_likes
        SET is_like = ?
        WHERE user_id = ? AND post_id = ?`,
		isLike, userID, postID)

	return err
}
