package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

/*
ログインをして、商品一覧ページに多めにアクセスする (画像読み込み含む)
サイトに負荷をかけているのに、商品を買わない嫌なユーザ
*/
func justLookingScenario(wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	score := 0
	resp := 200
	var c []*http.Cookie

	_, email, password := getUserInfo(0)
	resp, c = postLogin(c, email, password)
	score = calcScore(score, resp)

	resp, c = getIndex(c, 0)
	score = calcScore(score, resp)

	for i := 0; i < 50; i++ {
		resp, c = getImage(c, i%5)
		score = calcScore(score, resp)
	}
	updateScore(score, wg, m, finishTime)
	score = 0

	resp, c = getIndex(c, getRand(50, 99))
	score = calcScore(score, resp)

	resp, c = getIndex(c, getRand(100, 149))
	score = calcScore(score, resp)

	for i := 0; i < 50; i++ {
		resp, c = getImage(c, i%5)
		score = calcScore(score, resp)
	}
	updateScore(score, wg, m, finishTime)
	score = 0

	resp, c = getIndex(c, getRand(150, 199))
	score = calcScore(score, resp)

	resp, c = getProduct(c, 0)
	score = calcScore(score, resp)

	resp, c = getProduct(c, 0)
	score = calcScore(score, resp)

	resp, c = getProduct(c, 0)
	score = calcScore(score, resp)

	resp, c = getLogout(c)
	score = calcScore(score, resp)

	updateScore(score, wg, m, finishTime)
}

/*
ログインしないで、ユーザページに多めにアクセスする
他人の購入履歴を見てニヤニヤしている、ネットストーカー
*/
func stalkerScenario(wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	score := 0
	resp := 200
	var c []*http.Cookie

	resp, c = getIndex(c, 0)
	score = calcScore(score, resp)

	// id:1234 よく商品を買うユーザ
	resp, c = getUserPage(c, 1234)
	score = calcScore(score, resp)

	resp, c = getUserPage(c, 0)
	score = calcScore(score, resp)

	resp, c = getUserPage(c, 0)
	score = calcScore(score, resp)

	resp, c = getUserPage(c, 0)
	score = calcScore(score, resp)

	updateScore(score, wg, m, finishTime)
}

/*
ひたすら商品を買って、コメントをする
近年経済成長し、品質の高い先進国の商品を買いたくなっている人
*/
func bakugaiScenario(wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	score := 0
	resp := 200
	var c []*http.Cookie

	// 1/3 の確率で id:1234 のユーザが爆買いする
	uID := 0
	if getRand(1, 2) == 1 {
		uID = 1234
	}

	_, email, password := getUserInfo(uID)
	resp, c = postLogin(c, email, password)
	score = calcScore(score, resp)

	resp, c = getIndex(c, getRand(100, 199))
	score = calcScore(score, resp)

	for i := 0; i < 20; i++ {
		resp, c = buyProduct(c, 0)
		score = calcScore(score, resp)
	}
	updateScore(score, wg, m, finishTime)
	score = 0

	for i := 0; i < 5; i++ {
		resp, c = sendComment(c, 0)
		score = calcScore(score, resp)
	}

	resp, c = getLogout(c)
	score = calcScore(score, resp)

	updateScore(score, wg, m, finishTime)
}

// 以下、スコア計算用
func updateScore(score int, wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	m.Lock()
	defer m.Unlock()
	totalScore = totalScore + score
	if time.Now().After(finishTime) {
		wg.Done()
		if finished == false {
			finished = true
			showScore()
			postScore()
		}
	}
}

func calcScore(score int, response int) int {
	if response == 200 {
		return score + 1
	} else if strings.Contains(strconv.Itoa(response), "4") {
		return score - 20
	} else {
		return score - 50
	}
}

func showScore() {
	log.Print("Benchmark Finish!")
	log.Print("Score: " + strconv.Itoa(totalScore))
	log.Print("Waiting for Stopping All Benchmarkers ...")
}

func postScore() {
	apiURL := os.Getenv("BENCH_PORTAL_APIGW_URL")
	teamName := os.Getenv("BENCH_TEAM_NAME")
	if apiURL == "" && teamName == "" {
		return
	}

	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Print("Failed to send score")
		log.Printf("Error loading location: %v\n", err)
		return
	}
	now := time.Now().In(location)
	timestamp := now.Format(time.RFC3339)

	// Define the data to be sent
	data := map[string]interface{}{
		"team":      teamName,
		"score":     totalScore,
		"timestamp": timestamp,
	}

	// Convert the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Print("Failed to send score")
		log.Printf("Error encoding JSON: %v\n", err)
		return
	}

	// Create the PUT request
	req, err := http.NewRequest("PUT", apiURL+"teams", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Print("Failed to send score")
		log.Printf("Error creating request: %v\n", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the http.DefaultClient
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Failed to send score")
		log.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	log.Print("Sent score to portal site")
}
