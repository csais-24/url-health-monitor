package healthcheck

import "time"

type StatusCode int

type HealthCheckResult struct {
	URL        string
	StatusCode *StatusCode
	Reachable  bool
	Healthy    bool
	Error      error
	Duration   time.Duration
}
