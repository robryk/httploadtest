package main

import "flag"
import "fmt"
import "net/http"
import "os"
import "time"
import "github.com/robryk/httploadtest/stats"

var url = flag.String("url", "", "Url to loadtest")
var satoriToken = flag.String("token", "", "Satori token")
var maxConcurrent = flag.Int("max_concurrent_reqs", 10, "Maximum number of concurrent requests")

var tokens chan struct{}

var httpTransport http.Transport = http.Transport{
	MaxIdleConnsPerHost: 20,
}

var freqCounter stats.FreqCounter = stats.NewFreqCounter(stats.PrintCollector{
	Output: os.Stdout,
	Name:   "QPS",
})

var errFreqCounter stats.FreqCounter = stats.NewFreqCounter(stats.PrintCollector{
	Output: os.Stdout,
	Name: "EPS",
})

var latencyCollector stats.ValueCollector = stats.PrintCollector{
	Output: os.Stdout,
	Name:   "Latency",
}

func singleTest() {
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		panic(err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "satori_token",
		Value: *satoriToken,
	})
	c := http.Client{
		Transport: &httpTransport,
	}
	startTime := time.Now()
	var resp *http.Response
	resp, err = c.Do(req)
	if err == nil {
		resp.Body.Close()
		latencyCollector.Collect(time.Since(startTime).Seconds())
		freqCounter.Trigger()
	} else {
		fmt.Printf("Error: %v\n", err)
		time.Sleep(time.Second)
		errFreqCounter.Trigger()
	}
	tokens <- struct{}{}
}

func main() {
	flag.Parse()
	tokens = make(chan struct{}, *maxConcurrent)
	for i := 0; i < *maxConcurrent; i++ {
		tokens <- struct{}{}
	}
	freqCounter.Start()
	errFreqCounter.Start()
	for _ = range tokens {
		go singleTest()
	}
}
