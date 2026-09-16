package healthcheck

import (
	"net/http"
	"sync"
	"time"
)

type HealthChecker struct {
	client *http.Client
}

func (hc *HealthChecker) Check(url string) HealthCheckResult {
	var result HealthCheckResult
	result.URL = url
	start := time.Now()
	resp, err := hc.client.Get(url)
	result.Duration = time.Since(start)
	if err != nil {
		result.Error = err
		return result
	}
	defer resp.Body.Close()
	result.Reachable = true
	statusCode := StatusCode(resp.StatusCode)
	result.StatusCode = &statusCode

	if *result.StatusCode >= 200 && *result.StatusCode <= 299 {
		result.Healthy = true
	}

	return result
}

func (hc *HealthChecker) doWork(wg *sync.WaitGroup, urls <-chan string, res chan<- HealthCheckResult) {
	defer wg.Done()
	for url := range urls {
		res <- hc.Check(url)
	}
}

func (hc *HealthChecker) CheckAll(urls []string) []HealthCheckResult {
	var results []HealthCheckResult
	var wg sync.WaitGroup
	jobs := make(chan string)
	resultsChannel := make(chan HealthCheckResult, 20)
	workers := 30

	for range workers {
		wg.Add(1)
		go hc.doWork(&wg, jobs, resultsChannel)
	}

	go func() {
		for _, url := range urls {
			jobs <- url
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(resultsChannel)
	}()

	for result := range resultsChannel {
		results = append(results, result)
	}

	return results
}

func NewHealthChecker(client *http.Client) *HealthChecker {
	return &HealthChecker{
		client: client,
	}
}
