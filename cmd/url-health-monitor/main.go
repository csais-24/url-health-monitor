package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/csais-24/url-health-monitor/internal/healthcheck"
)

func main() {

	var url string
	var urlsFilePath string
	var results []healthcheck.HealthCheckResult
	flag.StringVar(&url, "url", "", "Individual URL to check health")
	flag.StringVar(&urlsFilePath, "url-list", "", "Path to URL file to check health")
	flag.Parse()

	if url != "" && urlsFilePath != "" {
		log.Fatal("You must select only 1 option. Either -url or -url-list")
	}

	if url == "" && urlsFilePath == "" {
		log.Fatalf("Usage: go run ./cmd/url-health-monitor/ -url <https://.example.com>\nOR\ngo run ./cmd/url-health-monitor/ -url-list <path-to-file.txt>")
	}

	client := &http.Client{Timeout: 5 * time.Second}
	hc := healthcheck.NewHealthChecker(client)

	if url != "" {
		results = hc.CheckAll([]string{url})
	} else {
		urls, err := parseUrlFile(urlsFilePath)
		if err != nil {
			log.Fatal(err)
		}
		if len(urls) == 0 {
			log.Fatal("no URLs were provided")
		}
		results = hc.CheckAll(urls)
	}

	for _, result := range results {
		printResult(result)
	}

	printSummary(results)

}

func parseUrlFile(urlsPath string) ([]string, error) {
	var urls []string
	b, err := os.ReadFile(urlsPath)
	if err != nil {
		return nil, err
	}
	rawUrls := strings.Split(string(b), "\n")
	for _, url := range rawUrls {
		trimmedUrl := strings.TrimSpace(url)
		if trimmedUrl == "" {
			continue
		}
		urls = append(urls, trimmedUrl)
	}

	return urls, nil
}

func getHealthStatus(healthy, reachable bool) string {
	if !reachable {
		return "UNREACHABLE"
	}

	if !healthy {
		return "UNHEALTHY"
	}

	return "HEALTHY"
}

func printResult(result healthcheck.HealthCheckResult) {
	status := getHealthStatus(result.Healthy, result.Reachable)
	if result.StatusCode == nil {
		fmt.Printf(
			"URL: %s\nStatus: %s\nError: %s\nDuration: %s\n\n",
			result.URL,
			status,
			result.Error,
			result.Duration,
		)
		return
	}

	fmt.Printf(
		"URL: %s\nStatus: %s\nHTTP Status: %d\nDuration: %s\n\n",
		result.URL,
		status,
		*result.StatusCode,
		result.Duration,
	)

}

func printSummary(results []healthcheck.HealthCheckResult) {
	total := len(results)
	healthy := 0
	unhealthy := 0
	unreachable := 0

	for _, result := range results {
		if !result.Reachable {
			unreachable++
		} else if result.Healthy {
			healthy++
		} else {
			unhealthy++
		}

	}

	fmt.Println("-------------------------------------------------")
	fmt.Println("Summary")
	fmt.Println("-------------------------------------------------")
	fmt.Printf("Total: \t\t\t%d\n", total)
	fmt.Printf("Healthy: \t\t%d\n", healthy)
	fmt.Printf("Unhealthy: \t\t%d\n", unhealthy)
	fmt.Printf("Unreachable: \t\t%d\n", unreachable)
}
