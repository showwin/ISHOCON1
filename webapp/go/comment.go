package main

// Comment Model
type Comment struct {
	ID        int
	ProductID int
	UserID    int
	Content   string
	CreatedAt string
}

func getComments(pid int) []Comment {
	rows, err := db.Query("SELECT * FROM comments WHERE product_id = ? ", pid)
	if err != nil {
		return nil
	}

	defer rows.Close()
	comments := []Comment{}
	for rows.Next() {
		c := Comment{}
		err = rows.Scan(&c.ID, &c.ProductID, &c.UserID, &c.Content, &c.CreatedAt)
		comments = append(comments, c)
	}

	return comments
}
