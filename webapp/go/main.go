package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	// load templates
	r.LoadHTMLGlob("templates/*")
	r.Use(static.Serve("/", static.LocalFile("public", true)))
	layout := "templates/layout.tmpl"

	// session store
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login", gin.H{
			"Message": "ECサイトで爆買いしよう！！！！",
		})
	})

	r.POST("/login", func(c *gin.Context) {
		email := c.PostForm("email")
		pass := c.PostForm("password")

		session := sessions.Default(c)
		uid, uPass := getUserPassword(email)
		if pass == uPass {
			// ログイン成功
			session.Set("uid", uid)
			session.Save()
		} else {
			// ログイン失敗
		}

		r.SetHTMLTemplate(template.Must(template.ParseFiles(layout, "templates/index.tmpl")))
		c.HTML(http.StatusOK, "base", gin.H{
			"title": "Main website",
		})
	})

	r.Run(":8080")
}

func db() *sql.DB {
	user := os.Getenv("ISHOCON1_DB_USER")
	pass := os.Getenv("ISHOCON1_DB_PASSWORD")
	dbname := "ishocon1"
	db, err := sql.Open("mysql", user+":"+pass+"@/"+dbname)
	if err != nil {
		panic(err.Error())
	}
	return db
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
