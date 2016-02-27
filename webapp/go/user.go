package main

import "github.com/gin-gonic/contrib/sessions"

// User model
type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	LastLogin string
}

func getUserPassword(email string) (int, string) {
	var uid int
	var pass string
	err := db().QueryRow("SELECT id, password FROM users WHERE email = ? LIMIT 1", email).Scan(&uid, &pass)
	if err != nil {
		panic(err.Error())
	}

	return uid, pass
}

func currentUser(session sessions.Session) User {
	id := session.Get("uid")
	u := User{}
	r := db().QueryRow("SELECT * FROM users WHERE id = ? LIMIT 1", id)
	err := r.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.LastLogin)
	if err != nil {
		panic(err.Error())
	}

	return u
}
