package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/csais-24/url-health-monitor/internal/healthcheck"
)

func main() {
	// For testing purpose
	var url string
	flag.StringVar(&url, "test-url", "", "Test URL")
	flag.Parse()
	if url == "" {
		log.Fatal("Usage: go run cmd/healthcheck/main.go <https://example.com")
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}
	hc := healthcheck.NewHealthChecker(client)
	hcResult := hc.Check(url)
	if hcResult.Error != nil {
		fmt.Printf("Url: %s\nError: %s\nDuration: %s\n", url, hcResult.Error, hcResult.Duration)
		return
	}
	fmt.Printf("Url: %s\nHealthy: %t\nStatus Code:%d\nReachable: %t\n", url, hcResult.Healthy, *hcResult.StatusCode, hcResult.Reachable)
}
