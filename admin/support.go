package main

import (
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func choice(s []string) string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(s))
	return s[i]
}

// ランダムにユーザー情報を取得
func getUserInfo(id int) (int, string, string) {
	if id == 0 {
		id = getRand(1, 5000)
	}
	var email, password string
	db, err := getDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.QueryRow("SELECT email, password FROM users WHERE id = ? LIMIT 1", id).Scan(&email, &password)
	if err != nil {
		panic(err.Error())
	}

	return id, email, password
}

// from から to までの値をランダムに取得
func getRand(from int, to int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(to+1-from) + from
}
