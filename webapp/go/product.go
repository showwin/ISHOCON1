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

func getProduct(pid int) Product {
	p := Product{}
	row := db().QueryRow("SELECT * FROM products WHERE id = ? LIMIT 1", pid)
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.ImagePath, &p.Price, &p.CreatedAt)
	if err != nil {
		panic(err.Error())
	}
	return p
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
