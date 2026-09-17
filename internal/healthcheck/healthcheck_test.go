package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestHealthChecker(t *testing.T) {
	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		ReachableExpected  bool
		HealthyExpected    bool
		Timeout            bool
		ExpectedError      bool
	}{
		{"HTTP Status OK Test", http.StatusOK, true, true, false, false},
		{"HTTP Status Internal Server Error", http.StatusInternalServerError, true, false, false, false},
		{"Timeout Test", 0, false, false, true, true},
	}

	client := &http.Client{Timeout: time.Millisecond * 500}
	hc := NewHealthChecker(client)
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.Timeout {
					time.Sleep(1 * time.Second)
					return
				}
				w.WriteHeader(tt.ExpectedStatusCode)
			}))
			res := hc.Check(srv.URL)
			if (res.Error != nil) != tt.ExpectedError {
				t.Errorf("Error should be nil got %s instead\n", res.Error)
				return
			}

			if res.StatusCode == nil {
				if !tt.ExpectedError {
					t.Errorf("Expected status code. Got nil instead")
					return
				}
				return
			}

			if int(*res.StatusCode) != tt.ExpectedStatusCode {
				t.Errorf("Expected status code %d -> Got %d instead\n", tt.ExpectedStatusCode, *res.StatusCode)
			}

			if res.Reachable != tt.ReachableExpected {
				t.Errorf("Expected Reachable to be %t. Got %t instead\n", tt.ReachableExpected, res.Reachable)
			}

			if res.Healthy != tt.HealthyExpected {
				t.Errorf("Expected Healthy to be %t. Got %t instead\n", tt.HealthyExpected, res.Healthy)
			}
			srv.Close()
		})

	}
}

func TestHealthChecker_CheckAll(t *testing.T) {
	testCasesLen := 10
	var urls []string
	var activeRequest atomic.Int32
	var maxConcurrentRequests atomic.Int32

	client := &http.Client{Timeout: 2 * time.Second}
	hc := NewHealthChecker(client)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activeRequest.Add(1)
		defer activeRequest.Add(-1)
		time.Sleep(1 * time.Millisecond)
		active := activeRequest.Load()

		for {
			max := maxConcurrentRequests.Load()
			if active <= max {
				break
			}

			if maxConcurrentRequests.CompareAndSwap(max, active) {
				break
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	for range testCasesLen {
		urls = append(urls, srv.URL)
	}

	results := hc.CheckAll(urls)
	if maxConcurrentRequests.Load() <= 1 {
		t.Errorf("maxConcurrentRequests = %d expected > 1\n", maxConcurrentRequests.Load())
	}

	if activeRequest.Load() != 0 {
		t.Error("activeRequest > 0 expected 0\n")
	}

	if len(results) != testCasesLen {
		t.Errorf("Expected result's len to be %d got %d instead\n", testCasesLen, len(results))
	}

	for _, result := range results {
		if result.URL == "" {
			t.Errorf("Expected url value got empty instead\n")
			continue
		}

		if result.StatusCode == nil {
			t.Errorf("Expected status code %d got <nil> instead", http.StatusOK)
			continue
		}

		if *result.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d got %d instead\n", http.StatusOK, *result.StatusCode)
		}

		if !result.Reachable {
			t.Errorf("Expected reachable to be true got false instead")
		}

		if !result.Healthy {
			t.Errorf("Expected healthy to be true got false instead")
		}

		if result.Error != nil {
			t.Errorf("Expected error to be nil got %s instead\n", result.Error)
		}
	}
}
