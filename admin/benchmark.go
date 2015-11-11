package main

import (
	"flag"
	"fmt"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func LoopJustLookingScenario(wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	for {
		JustLookingScenario(wg, m, finishTime)
	}
}

func LoopStalkerScenario(wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	for {
		StalkerScenario(wg, m, finishTime)
	}
}

func LoopBakugaiScenario(wg *sync.WaitGroup, m *sync.Mutex, finishTime time.Time) {
	for {
		BakugaiScenario(wg, m, finishTime)
	}
}

func StartBenchmark(workload int) {
	GetInitialize()
	ShowLog("Benchmark Start!  Workload: " + strconv.Itoa(workload))
	finishTime := time.Now().Add(1 * time.Minute)
	ValidateInitialize()
	wg := new(sync.WaitGroup)
	m := new(sync.Mutex)
	for i := 0; i < workload; i++ {
		wg.Add(1)
		if i%3 == 0 {
			go LoopJustLookingScenario(wg, m, finishTime)
		} else if i%3 == 1 {
			go LoopStalkerScenario(wg, m, finishTime)
		} else {
			go LoopBakugaiScenario(wg, m, finishTime)
		}
	}
	wg.Wait()
}

var host = "http://127.0.0.1"
var TotalScore = 9
var Finished = false

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

	StartBenchmark(*workload)
}
