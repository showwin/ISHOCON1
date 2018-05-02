package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

// 初期化確認
func validateInitialize() {
	validateIndex(10, false)
	validateProducts(false)
	validateUsers(1500, false)
	id, email, password := getUserInfo(0)
	var c []*http.Cookie
	_, c = postLogin(c, email, password)
	buyProduct(c, 10000)
	validateUsers(id, true)
	sendComment(c, 10000)
	validateIndex(0, true)
}

func validateIndex(page int, loggedIn bool) {
	var flg, flg1, flg2, flg3 bool
	var flg50, flgOrder, flgReview, flgPrdExp bool
	doc, err := goquery.NewDocument(host + "/?page=" + strconv.Itoa(page))
	if err != nil {
		log.Print("Cannot GET /index")
		os.Exit(1)
	}

	// 商品が50個あることの確認
	flg50 = doc.Find(".row").Children().Size() == 50

	// 商品が id DESC 順に並んでいることの確認
	doc.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// ついでにログインボタンの確認
		if i == 1 {
			ref, _ := s.Attr("href")
			flg1 = ref == "/login"
		}
		if i == 2 {
			ref, _ := s.Attr("href")
			flg2 = ref == "/products/"+strconv.Itoa(10000-page*50)
		}
		if i == 10 {
			ref, _ := s.Attr("href")
			flg3 = ref == "/products/"+strconv.Itoa(10000-page*50-4)
			return false
		}
		return true
	})
	flgOrder = flg1 && flg2 && flg3

	// レビューの件数が正しいことの確認
	doc.Find("h4").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 2 {
			str := s.Text()
			if loggedIn {
				flg1 = str == "21件のレビュー"
			} else {
				flg1 = str == "20件のレビュー"
			}
		}
		if i == 11 {
			str := s.Text()
			flg2 = str == "20件のレビュー"
			return false
		}
		return true
	})
	flgReview = flg1 && flg2

	// 商品のDOMの確認
	flg1 = doc.Find(".panel-default").First().Children().Size() == 2
	flg2 = doc.Find(".panel-body").First().Children().Size() == 7
	flg3 = doc.Find(".col-md-4").First().Find(".panel-body ul").Children().Size() == 5
	flgPrdExp = flg1 && flg2 && flg3

	// イメージパスの確認
	doc.Find("img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 0 || i == 6 || i == 12 || i == 18 || i == 24 {
			src, _ := s.Attr("src")
			flg1 = src == "/images/image"+strconv.Itoa((99-i)%5)+".jpg"
			flgPrdExp = flgPrdExp && flg1
		}
		if i == 24 {
			return false
		}
		return true
	})

	flg = flg50 && flgOrder && flgReview && flgPrdExp
	// 全体の確認
	if flg == false {
		log.Print("Invalid Content or DOM at GET /index")
		if flg50 == false {
			log.Print("商品が50個表示されていません")
		}
		if flgOrder == false {
			log.Print("商品が正しい順で並んでいません")
		}
		if flgReview == false {
			log.Print("レビューの件数が正しくありません")
		}
		if flgPrdExp == false {
			log.Print("商品説明部分が正しくありません")
		}
		os.Exit(1)
	}
}

func validateProducts(loggedIn bool) {
	var flg bool
	var flgImage, flgPrdExp bool
	doc, err := goquery.NewDocument(host + "/products/1500")
	if err != nil {
		log.Print("Cannot GET /products/:id")
		os.Exit(1)
	}
	// 画像パスの確認
	doc.Find("img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		src, _ := s.Attr("src")
		flgImage = src == "/images/image4.jpg"
		return false
	})

	// DOMの構造確認
	flgPrdExp = doc.Find(".row div.jumbotron").Children().Size() == 5

	// 商品説明確認
	doc.Find(".row div.jumbotron p").Each(func(i int, s *goquery.Selection) {
		if i == 1 {
			str := s.Text()
			flgPrdExp = strings.Contains(str, "1499") && flgPrdExp
		}
	})

	// 購入済み文章の確認(なし)
	flgPrdExp = doc.Find(".jumbotron div.container").Children().Size() == 1 && flgPrdExp

	flg = flgImage && flgPrdExp
	// 全体の確認
	if flg == false {
		log.Print("Invalid Content or DOM at GET /products/:id")
		if flgImage == false {
			log.Print("商品の画像が正しくありません")
		}
		if flgPrdExp == false {
			log.Print("商品説明部分が正しくありません")
		}
		os.Exit(1)
	}
}

func validateUsers(id int, loggedIn bool) {
	var flg bool
	var flg30, flgDOM, flgTotal, flgTime bool
	doc, err := goquery.NewDocument(host + "/users/" + strconv.Itoa(id))
	if err != nil {
		log.Print("Cannot GET /users/:id")
		os.Exit(1)
	}

	// 履歴が30個あることの確認
	flg30 = doc.Find(".row").Children().Size() == 30

	// DOMの確認
	flgDOM = doc.Find(".panel-default").First().Children().Size() == 2
	flgDOM = doc.Find(".panel-body").First().Children().Size() == 7 && flgDOM

	// 合計金額の確認
	sum := getTotalPay(id)
	doc.Find(".container h4").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		str := s.Text()
		flgTotal = str == "合計金額: "+sum+"円"
		return false
	})

	flgTime = true
	if loggedIn {
		// 一番最初に、最後に買った商品が出ていることの確認
		doc.Find(".panel-heading a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			str, _ := s.Attr("href")
			flgTime = strings.Contains(str, "10000") && flgTime
			return false
		})
		// 商品の購入時間が最近10秒以内であることの確認
		doc.Find(".panel-body p").Each(func(i int, s *goquery.Selection) {
			if i == 2 {
				str := s.Text()
				timeformat := "2006-01-02 15:04:05 -0700"
				createdAt, _ := time.Parse(timeformat, str+" +0900")
				flgTime = time.Now().Before(createdAt.Add(10*time.Second)) && flgTime
			}
		})
	}

	// 全体の確認
	flg = flg30 && flgDOM && flgTotal && flgTime
	if flg == false {
		log.Print("Invalid Content or DOM at GET /users/:id")
		if flg30 == false {
			log.Print("購入履歴の数が正しくありません")
		}
		if flgDOM == false {
			log.Print("UserページのDOMが正しくありません")
		}
		if flgTotal == false {
			log.Print("購入金額の合計が正しくありません")
		}
		if flgTime == false {
			log.Print("最近の購入履歴が正しくありません")
		}
		os.Exit(1)
	}
}

func getTotalPay(userID int) string {
	user := os.Getenv("ISHOCON1_DB_USER")
	pass := os.Getenv("ISHOCON1_DB_PASSWORD")
	dbname := os.Getenv("ISHOCON1_DB_NAME")
	db, err := sql.Open("mysql", user+":"+pass+"@/"+dbname)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query := `
    SELECT SUM(p.price) as total_pay
    FROM histories as h
    INNER JOIN products as p
    ON p.id = h.product_id
    WHERE h.user_id = ?`
	var totalPay string
	err = db.QueryRow(query, userID).Scan(&totalPay)
	if err != nil {
		panic(err.Error())
	}

	return totalPay
}
