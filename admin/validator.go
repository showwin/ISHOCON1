package main

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

// 初期化確認
func ValidateInitialize() {
	ValidateIndex(10, false)
	ValidateProducts(false)
	ValidateUsers(1500, false)
	id, email, password := GetUserInfo(0)
	var c []*http.Cookie
	_, c = PostLogin(c, email, password)
	BuyProduct(c, 10000)
	ValidateUsers(id, true)
	SendComment(c, 10000)
	ValidateIndex(0, true)
}

func ValidateIndex(page int, loggedIn bool) {
	var flg, flg1, flg2, flg3 bool
	doc, err := goquery.NewDocument(host + "/?page=" + strconv.Itoa(page))
	if err != nil {
		ShowLog("Cannot GET /index")
		os.Exit(1)
	}

	// 商品が50個あることの確認
	flg = doc.Find(".row").Children().Size() == 50

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
	flg = flg && flg1 && flg2 && flg3

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
	flg = flg && flg1 && flg2

	// 商品のDOMの確認
	flg1 = doc.Find(".panel-default").First().Children().Size() == 2
	flg2 = doc.Find(".panel-body").First().Children().Size() == 7
	flg3 = doc.Find(".col-md-4").First().Find(".panel-body ul").Children().Size() == 5
	flg = flg && flg1 && flg2 && flg3

	// イメージパスの確認
	doc.Find("img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 0 || i == 6 || i == 12 || i == 18 || i == 24 {
			src, _ := s.Attr("src")
			flg1 = src == "/images/image"+strconv.Itoa((99-i)%5)+".jpg"
			flg = flg && flg1
		}
		if i == 24 {
			return false
		}
		return true
	})

	// 全体の確認
	if flg == false {
		ShowLog("Invalid Content or DOM at GET /index")
		os.Exit(1)
	}
}

func ValidateProducts(loggedIn bool) {
	var flg bool
	doc, err := goquery.NewDocument(host + "/products/1500")
	if err != nil {
		ShowLog("Cannot GET /products/:id")
		os.Exit(1)
	}
	// 画像パスの確認
	doc.Find("img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		src, _ := s.Attr("src")
		flg = src == "/images/image4.jpg"
		return false
	})

	// DOMの構造確認
	flg = doc.Find(".row div.jumbotron").Children().Size() == 5 && flg

	// 商品説明確認
	doc.Find(".row div.jumbotron p").Each(func(i int, s *goquery.Selection) {
		if i == 1 {
			str := s.Text()
			flg = strings.Contains(str, "1499") && flg
		}
	})

	// 購入済み文章の確認(なし)
	flg = doc.Find(".jumbotron div.container").Children().Size() == 1 && flg

	// 全体の確認
	if flg == false {
		ShowLog("Invalid Content or DOM at GET /products/:id")
		os.Exit(1)
	}
}

func ValidateUsers(id int, loggedIn bool) {
	var flg bool
	doc, err := goquery.NewDocument(host + "/users/" + strconv.Itoa(id))
	if err != nil {
		ShowLog("Cannot GET /users/:id")
		os.Exit(1)
	}

	// 履歴が30個あることの確認
	flg = doc.Find(".row").Children().Size() == 30

	// DOMの確認
	flg = doc.Find(".panel-default").First().Children().Size() == 2 && flg
	flg = doc.Find(".panel-body").First().Children().Size() == 7 && flg

	// 合計金額の確認
	sum := GetTotalPay(id)
	doc.Find(".container h4").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		str := s.Text()
		flg = str == "合計金額: "+sum+"円" && flg
		return false
	})

	if loggedIn {
		// 一番最初に、最後に買った商品が出ていることの確認
		doc.Find(".panel-heading a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			str, _ := s.Attr("href")
			flg = strings.Contains(str, "10000") && flg
			return false
		})
		// 商品の購入時間が最近10秒以内であることの確認
		doc.Find(".panel-body p").Each(func(i int, s *goquery.Selection) {
			if i == 2 {
				str := s.Text()
				timeformat := "2006-01-02 15:04:05 -0700"
				created_at, _ := time.Parse(timeformat, str+" +0900")
				flg = time.Now().Before(created_at.Add(10*time.Second)) && flg
			}
		})
	}

	// 全体の確認
	if flg == false {
		ShowLog("Invalid Content or DOM at GET /users/:id")
		os.Exit(1)
	}
}

func GetTotalPay(user_id int) string {
	db, err := sql.Open("mysql", "ishocon:ishocon@/ishocon1")
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
	var total_pay string
	err = db.QueryRow(query, user_id).Scan(&total_pay)
	if err != nil {
		panic(err.Error())
	}

	return total_pay
}
