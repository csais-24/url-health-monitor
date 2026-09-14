package healthcheck

import (
	"net/http"
	"time"
)

type HealthChecker struct {
	Client *http.Client
}

func (hc *HealthChecker) Check(url string) HealthCheckResult {
	var result HealthCheckResult
	result.URL = url
	start := time.Now()
	resp, err := hc.Client.Get(url)
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
