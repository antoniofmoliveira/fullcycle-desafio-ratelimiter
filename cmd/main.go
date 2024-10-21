package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antoniofmoliveira/fullcycle-desafio-ratelimiter/internal/limiter"
)

// main starts a simple HTTP server that listens on the given port and responds to
// GET requests to /hello with "Hello, World!". It also implements rate limiting
// using a token bucket algorithm. The server can also simulate errors and slow
// requests.
func main() {

	port := flag.Int("port", 8080, "Port to listen on")
	qtTokens := flag.Int("qt-tokens", 10000, "Number of tokens in the token bucket. 0 to disable rate limiting")
	timeFrameSeconds := flag.Int("time-frame-seconds", 5, "Time frame in seconds. 0 to disable rate limiting")
	simulateErrors := flag.Bool("simulate-errors", true, "Simulate errors")
	simulateSlowRequests := flag.Bool("simulate-slow-requests", true, "Simulate slow requests")
	seedForSimulateErrors := flag.Int("seed-for-simulate-errors", 100, "Seed for simulate errors. 100 equals 1 in 100 requests randomly fail")
	seedForSimulateSlowRequests := flag.Int("seed-for-simulate-slow-requests", 100, "Seed for simulate slow requests (microseconds)")
	flag.Parse()

	server := &http.Server{Addr: fmt.Sprintf(":%d", *port)}

	var rateLimiter *limiter.TokenBucket
	if *qtTokens > 0 && *timeFrameSeconds > 0 {
		rateLimiter = limiter.NewTokenBucket(*qtTokens, time.Duration(*timeFrameSeconds)*time.Second)
	}

	handleHelloWord := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if *simulateErrors && rand.Intn(*seedForSimulateErrors) == 0 {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if *simulateSlowRequests {
			time.Sleep(time.Duration(rand.Intn(*seedForSimulateSlowRequests)) * time.Microsecond)
		}

		fmt.Fprintf(w, "Hello, World!")
	})

	if rateLimiter == nil {
		http.Handle("/hello", handleHelloWord)
	} else {
		http.Handle("/hello", rateLimiter.Middleware(handleHelloWord))
	}

	go func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v\n", err)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-termChan
	log.Println("server: shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not shutdown the server: %v\n", err)
	}
	fmt.Println("Server stopped")
	os.Exit(0)
}
