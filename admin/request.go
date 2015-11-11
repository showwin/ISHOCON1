package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 初期化(60秒以内)
func GetInitialize() {
	ShowLog("Start GET /initialize")
	finishTime := time.Now().Add(1 * time.Minute)
	HttpRequest("GET", "/initialize", nil, nil)
	if time.Now().Sub(finishTime) > 0 {
		ShowLog("Timeover at GET /initialize")
		os.Exit(1)
	}
}

// 基本アクセス
func GetIndex(c []*http.Cookie, page int) (int, []*http.Cookie) {
	return HttpRequest("GET", "/?page="+strconv.Itoa(page), nil, c)
}

func GetImage(c []*http.Cookie, id int) (int, []*http.Cookie) {
	return HttpRequest("GET", "/images/image"+strconv.Itoa(id)+".jpg", nil, c)
}

func GetProduct(c []*http.Cookie, id int) (int, []*http.Cookie) {
	if id == 0 {
		id = GetRand(1, 10000)
	}
	return HttpRequest("GET", "/products/"+strconv.Itoa(id), nil, c)
}

func GetUserPage(c []*http.Cookie, id int) (int, []*http.Cookie) {
	if id == 0 {
		id = GetRand(1, 5000)
	}
	return HttpRequest("GET", "/users/"+strconv.Itoa(id), nil, c)
}

func PostLogin(c []*http.Cookie, email string, password string) (int, []*http.Cookie) {
	v := url.Values{}
	v.Add("email", email)
	v.Add("password", password)
	return HttpRequest("POST", "/login", v, c)
}

func GetLogout(c []*http.Cookie) (int, []*http.Cookie) {
	return HttpRequest("GET", "/logout", nil, c)
}

func BuyProduct(c []*http.Cookie, productId int) (int, []*http.Cookie) {
	if productId == 0 {
		productId = GetRand(1, 10000)
	}
	return HttpRequest("POST", "/products/buy/"+strconv.Itoa(productId), nil, c)
}

func SendComment(c []*http.Cookie, productId int) (int, []*http.Cookie) {
	if productId == 0 {
		productId = GetRand(1, 10000)
	}
	v := url.Values{}
	opt := []string{"爆買いしてよかった。", "二度と買わない。", "友達にも勧めます。"}
	v.Add("content", strings.Repeat("この商品は"+choice(opt), 5))
	return HttpRequest("POST", "/comments/"+strconv.Itoa(productId), v, c)
}

func HttpRequest(method string, path string, params url.Values, cookies []*http.Cookie) (int, []*http.Cookie) {
	req, _ := http.NewRequest(method, host+path, strings.NewReader(params.Encode()))
	jar, _ := cookiejar.New(nil)
	CookieURL, _ := url.Parse(host + path)
	jar.SetCookies(CookieURL, cookies)
	client := http.Client{Jar: jar}

	resp, err := client.Do(req)
	if err != nil {
		return 500, cookies
	}
	defer resp.Body.Close()

	return resp.StatusCode, jar.Cookies(CookieURL)
}

func ShowScore() {
	ShowLog("Benchmark Finish!")
	ShowLog("Score: " + strconv.Itoa(TotalScore))
	ShowLog("Waiting for Stopping All Benchmarkers ...")
}

func ShowLog(str string) {
	fmt.Println(time.Now().Format("15:04:05") + "  " + str)
}
