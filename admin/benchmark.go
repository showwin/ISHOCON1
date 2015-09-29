package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"net/http"
  "net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// 初期化(60秒以内)
func GetInitialize() {
	ShowLog("Start GET /initialize")
	finishTime := time.Now().Add(1 * time.Minute)
	goquery.NewDocument(host + "/initialize")
	if time.Now().Sub(finishTime) > 0 {
		ShowLog("Timeover at GET /initialize")
		os.Exit(1)
	}
}

// 初期化確認
func ValidateInitialize() {
	ValidateIndex(10, false)
	ValidateProducts(false)
	ValidateUsers(1500, false)
	id, email, password := GetUserInfo()
	PostLogin(email, password)
	BuyProduct(10000)
  ValidateUsers(id, true)
  SendComment(10000)
  ValidateIndex(0, true)
}

func ValidateIndex(page int, loggedIn bool) {
	var flg, flg1, flg2, flg3 bool
	doc, _ := goquery.NewDocument(host + "/?page="+strconv.Itoa(page))

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
			flg2 = ref == "/products/" + strconv.Itoa(10000-page*50)
		}
		if i == 10 {
			ref, _ := s.Attr("href")
			flg3 = ref == "/products/" + strconv.Itoa(10000-page*50-4)
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
        flg1 = str == "レビュー(21件)"
      } else {
        flg1 = str == "レビュー(20件)"
      }
		}
		if i == 11 {
			str := s.Text()
			flg2 = str == "レビュー(20件)"
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

  // コメントの反映確認
  if loggedIn {
    doc.Find(".panel-body ul li").EachWithBreak(func(_ int, s *goquery.Selection) bool {
      str := s.Text()
      flg = strings.Contains(str, "この商品") && flg
      return false
    })
  }

	// 全体の確認
	if flg == false {
		ShowLog("Invalid Content or DOM at GET /index")
		os.Exit(1)
	}
}

func ValidateProducts(loggedIn bool) {
	var flg bool
	doc, _ := goquery.NewDocument(host + "/products/1500")

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
	doc, _ := goquery.NewDocument(host + "/users/"+strconv.Itoa(id))

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
        created_at, _ := time.Parse(timeformat, str)
        flg = time.Now().Before(created_at.Add(10 * time.Second)) && flg
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

// 基本アクセス
func GetIndex(page int) int {
  return HttpRequest("GET", "/?page="+strconv.Itoa(page), nil)
}

func GetImage(id int) int {
  return HttpRequest("GET", "/images/image"+strconv.Itoa(id)+".jpg", nil)
}

func GetProduct(id int) int {
  if id == 0 {
    id = GetRand(1, 10000)
  }
  return HttpRequest("GET", "/products/"+strconv.Itoa(id), nil)
}

func PostLogin(email string, password string) int {
	v := url.Values{}
	v.Add("email", email)
	v.Add("password", password)
  return HttpRequest("POST", "/login", v)
}

func GetLogout() int {
  return HttpRequest("GET", "/logout", nil)
}

func BuyProduct(productId int) int {
  if productId == 0 {
    productId = GetRand(1, 10000)
  }
  return HttpRequest("POST", "/products/buy/"+strconv.Itoa(productId), nil)
}

func SendComment(productId int) int {
  if productId == 0 {
    productId = GetRand(1, 10000)
  }
  v := url.Values{}
  c := []string{"爆買いしてよかった。", "二度と買わない。", "友達にも勧めます。"}
	v.Add("content", strings.Repeat("この商品は"+choice(c), 5))
  return HttpRequest("POST", "/comments/"+strconv.Itoa(productId), v)
}

func choice(s []string) string {
  rand.Seed(time.Now().UnixNano())
  i := rand.Intn(len(s))
  return s[i]
}

// ランダムにユーザー情報を取得
func GetUserInfo() (int, string, string) {
	id := GetRand(1, 5000)
	var email, password string
	db, err := sql.Open("mysql", "ishocon:ishocon@/ishocon1")
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
func GetRand(from int, to int) int {
  rand.Seed(time.Now().UnixNano())
	return rand.Intn(to+1-from) + from
}

func LoopAccess() int {
  s := 0
  // ログイン状態
  _, email, password := GetUserInfo()
	s = CalcScore(s, PostLogin(email, password))
  s = CalcScore(s, GetIndex(0))
  s = CalcScore(s, GetProduct(0))
  s = CalcScore(s, SendComment(0))
  for i := 0; i < 10; i++ {
    s = CalcScore(s, GetImage(0))
    s = CalcScore(s, GetImage(1))
    s = CalcScore(s, GetImage(2))
    s = CalcScore(s, GetImage(3))
    s = CalcScore(s, GetImage(4))
    s = CalcScore(s, BuyProduct(0))
  }
  //ログアウト状態
  s = CalcScore(s, GetLogout())
  s = CalcScore(s, GetIndex(GetRand(150, 199)))
	s = CalcScore(s, GetProduct(0))
  return s
}

func CalcScore(score int, response int) int {
  if response == 200 {
    return score + 1
  } else if strings.Contains(strconv.Itoa(response), "4") {
    return score - 20
  } else {
    return score - 50
  }
}

func HttpRequest(method string, path string, params url.Values) int {
  req, _ := http.NewRequest(method, host+path, strings.NewReader(params.Encode()))
  jar, _ := cookiejar.New(nil)
  CookieURL, _ := url.Parse(host+path)
  jar.SetCookies(CookieURL, cookies)
  client := http.Client{Jar: jar}

	resp, err := client.Do(req)
	if err != nil {
		return 500
	}
  defer resp.Body.Close()

  cookies = jar.Cookies(CookieURL)
  return resp.StatusCode
}

func StartBenchmark() {
	GetInitialize()
	ShowLog("Benchmark Start!")
	finishTime := time.Now().Add(1 * time.Minute)
	ValidateInitialize()
  score := 9
	for {
		if time.Now().After(finishTime) {
			break
		}
		score = score + LoopAccess()
	}
	ShowLog("Benchmark Finish!")
	ShowLog("Score: " + strconv.Itoa(score))
}

func ShowLog(str string) {
	fmt.Println(time.Now().Format("15:04:05") + "  " + str)
}

const host = "http://localhost:8080"
var cookies []*http.Cookie

func main() {
	StartBenchmark()
}
