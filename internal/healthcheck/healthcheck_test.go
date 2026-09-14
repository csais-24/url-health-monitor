package healthcheck

import (
	"net/http"
	"net/http/httptest"
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
	hc := HealthChecker{Client: client}
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
