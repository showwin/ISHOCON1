package main

import "log"

// Product Model
type Product struct {
	ID          int
	Name        string
	Description string
	ImagePath   string
	Price       int
	CreatedAt   string
}

// ProductWithComments Model
type ProductWithComments struct {
	ID           int
	Name         string
	Description  string
	ImagePath    string
	Price        int
	CreatedAt    string
	CommentCount int
	Comments     []CommentWriter
}

// CommentWriter Model
type CommentWriter struct {
	Content string
	Writer  string
}

func getProduct(pid int) Product {
	p := Product{}
	row := db().QueryRow("SELECT * FROM products WHERE id = ? LIMIT 1", pid)
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.ImagePath, &p.Price, &p.CreatedAt)
	if err != nil {
		panic(err.Error())
	}

	return p
}

func getProductsWithCommentsAt(page int) []ProductWithComments {
	// select 50 products with offset page*50
	products := []ProductWithComments{}
	rows, err := db().Query("SELECT * FROM products ORDER BY id DESC LIMIT 50 OFFSET ?", page*50)
	if err != nil {
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		p := ProductWithComments{}
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImagePath, &p.Price, &p.CreatedAt)

		// select comment count for the product
		var cnt int
		cnterr := db().QueryRow("SELECT count(*) as count FROM comments WHERE product_id = ?", p.ID).Scan(&cnt)
		if cnterr != nil {
			cnt = 0
		}
		p.CommentCount = cnt

		if cnt > 0 {
			// select 5 comments and its writer for the product
			var cWriters []CommentWriter

			subrows, suberr := db().Query("SELECT * FROM comments as c INNER JOIN users as u "+
				"ON c.user_id = u.id WHERE c.product_id = ? ORDER BY c.created_at DESC LIMIT 5", p.ID)
			if suberr != nil {
				subrows = nil
			}

			defer subrows.Close()
			for subrows.Next() {
				var i int
				var s string
				var cw CommentWriter
				subrows.Scan(&i, &i, &i, &cw.Content, &s, &i, &cw.Writer, &s, &s, &s)
				cWriters = append(cWriters, cw)
			}

			p.Comments = cWriters
		}

		products = append(products, p)
	}

	return products
}

func (p *Product) isBought(uid int) bool {
	var count int
	log.Print(uid)
	log.Print(p.ID)
	err := db().QueryRow(
		"SELECT count(*) as count FROM histories WHERE product_id = ? AND user_id = ?",
		p.ID, uid,
	).Scan(&count)
	if err != nil {
		panic(err.Error())
	}

	return count > 0
}
