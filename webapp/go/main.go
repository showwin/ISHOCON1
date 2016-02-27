package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	// load templates
	r.Use(static.Serve("/", static.LocalFile("public", true)))
	layout := "templates/layout.tmpl"

	// session store
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		tmpl, _ := template.ParseFiles("templates/login.tmpl")
		r.SetHTMLTemplate(tmpl)
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

		tmpl := template.Must(template.ParseFiles(layout))
		tmpl.ParseFiles("templates/index.tmpl")
		r.SetHTMLTemplate(tmpl)
		c.HTML(http.StatusOK, "base", gin.H{
			"title": "Main website",
		})
	})

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		tmpl, _ := template.ParseFiles("templates/login.tmpl")
		r.SetHTMLTemplate(tmpl)
		c.Redirect(http.StatusFound, "/login")
	})

	r.GET("/products/:productId", func(c *gin.Context) {
		pid, _ := strconv.Atoi(c.Param("productId"))
		product := getProduct(pid)
		comments := getComments(pid)

		cUser := currentUser(sessions.Default(c))
		bought := product.isBought(cUser.ID)

		r.SetHTMLTemplate(template.Must(template.ParseFiles(layout, "templates/product.tmpl")))
		c.HTML(http.StatusOK, "base", gin.H{
			"CurrentUser":   cUser,
			"Product":       product,
			"Comments":      comments,
			"AlreadyBought": bought,
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
