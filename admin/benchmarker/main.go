package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func loopJustLookingScenario(wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	for {
		justLookingScenario(wg, m, finishTime)
	}
}

func loopStalkerScenario(wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	for {
		stalkerScenario(wg, m, finishTime)
	}
}

func loopBakugaiScenario(wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	for {
		bakugaiScenario(wg, m, finishTime)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getDB() (*sql.DB, error) {
	user := getEnv("ISHOCON1_DB_USER", "ishocon")
	pass := getEnv("ISHOCON1_DB_PASSWORD", "ishocon")
	dbname := getEnv("ISHOCON1_DB_NAME", "ishocon1_bench")
	db, err := sql.Open("mysql", user+":"+pass+"@/"+dbname)
	return db, err
}

func startBenchmark(workload int) {
	getInitialize()
	log.Print("Benchmark Start!  Workload: " + strconv.Itoa(workload))
	finishTime := time.Now().Add(1 * time.Minute)
	validateInitialize()
	wg := new(sync.WaitGroup)
	m := new(sync.Mutex)
	for i := 0; i < workload; i++ {
		wg.Add(1)
		if i%3 == 0 {
			go loopJustLookingScenario(wg, m, finishTime)
		} else if i%3 == 1 {
			go loopStalkerScenario(wg, m, finishTime)
		} else {
			go loopBakugaiScenario(wg, m, finishTime)
		}
	}
	wg.Wait()
}

var host = "http://127.0.0.1"
var totalScore = 9
var finished = false

func main() {
	flag.Usage = func() {
		fmt.Println(`Usage: ./benchmark [option]
Options:
  --workload N	run benchmark with n workloads (default: 3)
  --ip IP	specify target ip (default: 127.0.0.1:80)`)
	}

	var (
		workload = flag.Int("workload", 3, "")
		ip       = flag.String("ip", "127.0.0.1", "")
	)
	flag.Parse()
	host = "http://" + *ip

	startBenchmark(*workload)
}
