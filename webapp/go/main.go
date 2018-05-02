package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"unicode/utf8"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// database setting
	user := os.Getenv("ISHOCON1_DB_USER")
	pass := os.Getenv("ISHOCON1_DB_PASSWORD")
	dbname := os.Getenv("ISHOCON1_DB_NAME")
	db, _ = sql.Open("mysql", user+":"+pass+"@/"+dbname)
	db.SetMaxIdleConns(5)

	r := gin.Default()
	// load templates
	r.Use(static.Serve("/css", static.LocalFile("public/css", true)))
	r.Use(static.Serve("/images", static.LocalFile("public/images", true)))
	layout := "templates/layout.tmpl"

	// session store
	store := sessions.NewCookieStore([]byte("mysession"))
	store.Options(sessions.Options{HttpOnly: true})
	r.Use(sessions.Sessions("showwin_happy", store))

	// GET /login
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

	// POST /login
	r.POST("/login", func(c *gin.Context) {
		email := c.PostForm("email")
		pass := c.PostForm("password")

		session := sessions.Default(c)
		user, result := authenticate(email, pass)
		if result {
			// 認証成功
			session.Set("uid", user.ID)
			session.Save()

			user.UpdateLastLogin()

			c.Redirect(http.StatusSeeOther, "/")
		} else {
			// 認証失敗
			tmpl, _ := template.ParseFiles("templates/login.tmpl")
			r.SetHTMLTemplate(tmpl)
			c.HTML(http.StatusOK, "login", gin.H{
				"Message": "ログインに失敗しました",
			})
		}
	})

	// GET /logout
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		tmpl, _ := template.ParseFiles("templates/login.tmpl")
		r.SetHTMLTemplate(tmpl)
		c.Redirect(http.StatusFound, "/login")
	})

	// GET /
	r.GET("/", func(c *gin.Context) {
		cUser := currentUser(sessions.Default(c))

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 0
		}
		products := getProductsWithCommentsAt(page)
		// shorten description and comment
		var sProducts []ProductWithComments
		for _, p := range products {
			if utf8.RuneCountInString(p.Description) > 70 {
				p.Description = string([]rune(p.Description)[:70]) + "…"
			}

			var newCW []CommentWriter
			for _, c := range p.Comments {
				if utf8.RuneCountInString(c.Content) > 25 {
					c.Content = string([]rune(c.Content)[:25]) + "…"
				}
				newCW = append(newCW, c)
			}
			p.Comments = newCW
			sProducts = append(sProducts, p)
		}

		r.SetHTMLTemplate(template.Must(template.ParseFiles(layout, "templates/index.tmpl")))
		c.HTML(http.StatusOK, "base", gin.H{
			"CurrentUser": cUser,
			"Products":    sProducts,
		})
	})

	// GET /users/:userId
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
			if utf8.RuneCountInString(p.Description) > 70 {
				p.Description = string([]rune(p.Description)[:70]) + "…"
			}
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

	// GET /products/:productId
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

	// POST /products/buy/:productId
	r.POST("/products/buy/:productId", func(c *gin.Context) {
		// need authenticated
		if notAuthenticated(sessions.Default(c)) {
			tmpl, _ := template.ParseFiles("templates/login.tmpl")
			r.SetHTMLTemplate(tmpl)
			c.HTML(http.StatusForbidden, "login", gin.H{
				"Message": "先にログインをしてください",
			})
		} else {
			// buy product
			cUser := currentUser(sessions.Default(c))
			cUser.BuyProduct(c.Param("productId"))

			// redirect to user page
			tmpl, _ := template.ParseFiles("templates/mypage.tmpl")
			r.SetHTMLTemplate(tmpl)
			c.Redirect(http.StatusFound, "/users/"+strconv.Itoa(cUser.ID))
		}
	})

	// POST /comments/:productId
	r.POST("/comments/:productId", func(c *gin.Context) {
		// need authenticated
		if notAuthenticated(sessions.Default(c)) {
			tmpl, _ := template.ParseFiles("templates/login.tmpl")
			r.SetHTMLTemplate(tmpl)
			c.HTML(http.StatusForbidden, "login", gin.H{
				"Message": "先にログインをしてください",
			})
		} else {
			// create comment
			cUser := currentUser(sessions.Default(c))
			cUser.CreateComment(c.Param("productId"), c.PostForm("content"))

			// redirect to user page
			tmpl, _ := template.ParseFiles("templates/mypage.tmpl")
			r.SetHTMLTemplate(tmpl)
			c.Redirect(http.StatusFound, "/users/"+strconv.Itoa(cUser.ID))
		}
	})

	// GET /initialize
	r.GET("/initialize", func(c *gin.Context) {
		db.Exec("DELETE FROM users WHERE id > 5000")
		db.Exec("DELETE FROM products WHERE id > 10000")
		db.Exec("DELETE FROM comments WHERE id > 200000")
		db.Exec("DELETE FROM histories WHERE id > 500000")

		c.String(http.StatusOK, "Finish")
	})

	r.Run(":8080")
}
