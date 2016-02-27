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
	store := sessions.NewCookieStore([]byte("showwin_happy"))
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
		uid, result := authenticate(email, pass)
		if result {
			// ログイン成功
			session.Set("uid", uid)
			session.Save()

			// TODO: redirect to /
			tmpl := template.Must(template.ParseFiles(layout))
			tmpl.ParseFiles("templates/index.tmpl")
			r.SetHTMLTemplate(tmpl)
			c.HTML(http.StatusOK, "base", gin.H{
				"title": "Main website",
			})
		} else {
			// ログイン失敗
			tmpl, _ := template.ParseFiles("templates/login.tmpl")
			r.SetHTMLTemplate(tmpl)
			c.HTML(http.StatusOK, "login", gin.H{
				"Message": "ログインに失敗しました",
			})
		}
	})

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		tmpl, _ := template.ParseFiles("templates/login.tmpl")
		r.SetHTMLTemplate(tmpl)
		c.Redirect(http.StatusFound, "/login")
	})

	r.GET("/", func(c *gin.Context) {

	})

	r.GET("/users/:userId", func(c *gin.Context) {
		cUser := currentUser(sessions.Default(c))

		uid, _ := strconv.Atoi(c.Param("userId"))
		user := getUser(uid)

		products := user.BuyingHistory()

		var totalPay int
		for _, p := range products {
			totalPay += p.Price
		}

		// shorten description
		var sdProducts []Product
		for _, p := range products {
			p.Description = string([]rune(p.Description)[:70])
			sdProducts = append(sdProducts, p)
		}

		r.SetHTMLTemplate(template.Must(template.ParseFiles(layout, "templates/mypage.tmpl")))
		c.HTML(http.StatusOK, "base", gin.H{
			"CurrentUser": cUser,
			"User":        user,
			"Products":    sdProducts,
			"TotalPay":    totalPay,
		})
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

	r.POST("/products/buy/:productId", func(c *gin.Context) {
		// need authenticated

	})

	r.POST("/comments/:product_id", func(c *gin.Context) {
		// need authenticated
	})

	r.GET("/initialize", func(c *gin.Context) {

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
