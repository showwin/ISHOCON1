package main

import (
	"time"

	"github.com/gin-gonic/contrib/sessions"
)

// User model
type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	LastLogin string
}

func authenticate(email string, password string) (User, bool) {
	var u User
	err := db.QueryRow("SELECT * FROM users WHERE email = ? LIMIT 1", email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.LastLogin)
	if err != nil {
		return u, false
	}
	result := u.Password == u.Password
	return u, result
}

func notAuthenticated(session sessions.Session) bool {
	uid := session.Get("uid")
	return !(uid.(int) > 0)
}

func getUser(uid int) User {
	u := User{}
	r := db.QueryRow("SELECT * FROM users WHERE id = ? LIMIT 1", uid)
	err := r.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.LastLogin)
	if err != nil {
		return u
	}

	return u
}

func currentUser(session sessions.Session) User {
	uid := session.Get("uid")
	u := User{}
	r := db.QueryRow("SELECT * FROM users WHERE id = ? LIMIT 1", uid)
	err := r.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.LastLogin)
	if err != nil {
		return u
	}

	return u
}

// BuyingHistory : products which user had bought
func (u *User) BuyingHistory() (products []Product) {
	rows, err := db.Query(
		"SELECT p.id, p.name, p.description, p.image_path, p.price, h.created_at "+
			"FROM histories as h "+
			"LEFT OUTER JOIN products as p "+
			"ON h.product_id = p.id "+
			"WHERE h.user_id = ? "+
			"ORDER BY h.id DESC", u.ID)
	if err != nil {
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		p := Product{}
		var cAt string
		fmt := "2006-01-02 15:04:05"
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImagePath, &p.Price, &cAt)
		tmp, _ := time.Parse(fmt, cAt)
		p.CreatedAt = (tmp.Add(9 * time.Hour)).Format(fmt)
		if err != nil {
			panic(err.Error())
		}
		products = append(products, p)
	}

	return
}

// BuyProduct : buy product
func (u *User) BuyProduct(pid string) {
	db.Exec(
		"INSERT INTO histories (product_id, user_id, created_at) VALUES (?, ?, ?)",
		pid, u.ID, time.Now())
}

// CreateComment : create comment to the product
func (u *User) CreateComment(pid string, content string) {
	db.Exec(
		"INSERT INTO comments (product_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
		pid, u.ID, content, time.Now())
}

func (u *User) UpdateLastLogin() {
	db.Exec("UPDATE users SET last_login = ? WHERE id = ?", time.Now(), u.ID)
}
