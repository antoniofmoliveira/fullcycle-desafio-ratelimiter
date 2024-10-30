package main

import (
	"net/http"
	"os"
	"testing"
	"time"
)

// TestServerRateLimiter tests the server's rate limiting functionality.
// It starts the server with a specific token bucket configuration and sends a series of requests.
// The test verifies that the server handles requests correctly according to the rate limit,
// allowing a certain number of successful requests and returning the expected status codes
// for requests that exceed the limit.
func TestServerRateLimiter(t *testing.T) {
	os.Args = []string{"cmd/main", "--qt-tokens", "3", "--time-frame-seconds", "2"}
	go main()                   // Start the server in a separate goroutine
	time.Sleep(1 * time.Second) // Give the server some time to start

	client := &http.Client{}
	requests := 5
	successfulRequests := 0

	for i := 0; i < requests; i++ {
		resp, err := client.Get("http://localhost:8080/hello")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			successfulRequests++
		} else if resp.StatusCode != http.StatusTooManyRequests {
			t.Errorf("Unexpected status code: got %v, want %v or %v", resp.StatusCode, http.StatusOK, http.StatusTooManyRequests)
		}
	}
	if successfulRequests != 3 {
		t.Errorf("Expected 3 successful requests, got %d", successfulRequests)
	}
	time.Sleep(3 * time.Second) // wait bucket reset

	requests = 5
	successfulRequests = 0

	for i := 0; i < requests; i++ {
		resp, err := client.Get("http://localhost:8080/hello")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			successfulRequests++
		} else if resp.StatusCode == http.StatusTooManyRequests {
			time.Sleep(3 * time.Second) // Wait bucket reset
		} else if resp.StatusCode != http.StatusTooManyRequests {
			t.Errorf("Unexpected status code: got %v, want %v or %v", resp.StatusCode, http.StatusOK, http.StatusTooManyRequests)

		}
	}
	if successfulRequests != 4 {
		t.Errorf("Expected 4 successful requests, got %d", successfulRequests)
	}

}
