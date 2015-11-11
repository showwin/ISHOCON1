package main

import (
	"net/http"
	"sync"
	"time"
  "strings"
  "strconv"
  "fmt"

	_ "github.com/go-sql-driver/mysql"
)

/*
  ログインをして、商品一覧ページに多めにアクセスする (画像読み込み含む)
  サイトに負荷をかけているのに、商品を買わない嫌なユーザ
*/
func JustLookingScenario(wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	score := 0
	resp := 200
	var c []*http.Cookie

	_, email, password := GetUserInfo(0)
	resp, c = PostLogin(c, email, password)
	score = CalcScore(score, resp)

	resp, c = GetIndex(c, 0)
	score = CalcScore(score, resp)

  for i := 0; i < 50; i++ {
		resp, c = GetImage(c, i%5)
		score = CalcScore(score, resp)
	}
  UpdateScore(score, wg, m, finishTime)
  score = 0

  resp, c = GetIndex(c, GetRand(50, 99))
	score = CalcScore(score, resp)

  resp, c = GetIndex(c, GetRand(100, 149))
	score = CalcScore(score, resp)

  for i := 0; i < 50; i++ {
		resp, c = GetImage(c, i%5)
		score = CalcScore(score, resp)
	}
  UpdateScore(score, wg, m, finishTime)
  score = 0

  resp, c = GetIndex(c, GetRand(150, 199))
	score = CalcScore(score, resp)

	resp, c = GetProduct(c, 0)
	score = CalcScore(score, resp)

  resp, c = GetProduct(c, 0)
	score = CalcScore(score, resp)

  resp, c = GetProduct(c, 0)
	score = CalcScore(score, resp)

  resp, c = GetLogout(c)
	score = CalcScore(score, resp)

	UpdateScore(score, wg, m, finishTime)
}

/*
  ログインしないで、ユーザページに多めにアクセスする
  他人の購入履歴を見てニヤニヤしている、ネットストーカー
*/
func StalkerScenario(wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	score := 0
	resp := 200
	var c []*http.Cookie

  resp, c = GetIndex(c, 0)
	score = CalcScore(score, resp)

  // id:1234 よく商品を買うユーザ
	resp, c = GetUserPage(c, 1234)
	score = CalcScore(score, resp)

  resp, c = GetUserPage(c, 0)
	score = CalcScore(score, resp)

  resp, c = GetUserPage(c, 0)
	score = CalcScore(score, resp)

  resp, c = GetUserPage(c, 0)
	score = CalcScore(score, resp)

	UpdateScore(score, wg, m, finishTime)
}

/*
  ひたすら商品を買って、コメントをする
  近年経済成長し、品質の高い先進国の商品を買いたくなっている人
*/
func BakugaiScenario(wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	score := 0
	resp := 200
	var c []*http.Cookie

  // 1/3 の確率で id:1234 のユーザが爆買いする
  uId := 0
  if GetRand(1, 2) == 1 {
    uId = 1234
  }

	_, email, password := GetUserInfo(uId)
	resp, c = PostLogin(c, email, password)
	score = CalcScore(score, resp)

  resp, c = GetIndex(c, GetRand(100, 199))
	score = CalcScore(score, resp)

  for i := 0; i < 20; i++ {
		resp, c = BuyProduct(c, 0)
		score = CalcScore(score, resp)
	}
  UpdateScore(score, wg, m, finishTime)
  score = 0

  for i := 0; i < 5; i++ {
		resp, c = SendComment(c, 0)
		score = CalcScore(score, resp)
	}

  resp, c = GetLogout(c)
	score = CalcScore(score, resp)

	UpdateScore(score, wg, m, finishTime)
}

// 以下、スコア計算用
func UpdateScore(score int, wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	m.Lock()
	defer m.Unlock()
	TotalScore = TotalScore + score
	if time.Now().After(finishTime) {
		wg.Done()
		if Finished == false {
			Finished = true
			ShowScore()
		}
	}
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
