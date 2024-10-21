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
	os.Args = []string{"cmd/main", "--qt-tokens", "3", "--time-frame-seconds", "5"}
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
}


// TestServerRateLimiterSurcharge tests the server's rate limiting behavior with surcharge.
// It starts the server with a specific token bucket configuration and sends a series of requests.
// The test verifies that after hitting the rate limit, requests wait for a surcharge period
// and then continue. It checks that the correct number of requests are successfully processed
// and that the expected status codes are returned.
func TestServerRateLimiterSurcharge(t *testing.T) {

	os.Args = []string{"cmd/main", "--qt-tokens", "3", "--time-frame-seconds", "5"}

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
		} else if resp.StatusCode == http.StatusTooManyRequests {
			time.Sleep(10 * time.Second) // Wait for 10 seconds before retrying the request
		} else if resp.StatusCode != http.StatusTooManyRequests {
			t.Errorf("Unexpected status code: got %v, want %v or %v", resp.StatusCode, http.StatusOK, http.StatusTooManyRequests)

		}
	}
	if successfulRequests != 4 {
		t.Errorf("Expected 4 successful requests, got %d", successfulRequests)
	}
}
