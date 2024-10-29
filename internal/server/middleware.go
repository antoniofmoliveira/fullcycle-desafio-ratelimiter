package server

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/antoniofmoliveira/fullcycle-desafio-ratelimiter/internal/config"
	"github.com/antoniofmoliveira/fullcycle-desafio-ratelimiter/internal/limiter"
)

// noSimulate is a middleware that simply forwards the request to the next handler
// in the chain, without performing any simulation. It is used when the
// simulation is disabled in the configuration.
func noSimulate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

// simulateSlowRequestsMiddleware returns a middleware that simulates slow requests
// based on the configured seed value. When the random condition is met, it delays
// the request by a duration defined in the configuration; otherwise, it passes
// the request to the next handler without delay.
func simulateSlowRequestsMiddleware(c *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c.SimulateSlowRequests && rand.Intn(c.SeedForSimulateSlowRequests) == 0 {
				time.Sleep(time.Duration(rand.Intn(c.SeedForSimulateSlowRequests)) * time.Microsecond)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// getSimulateSlowRequestsMiddleware returns a middleware that simulates slow requests
// based on the configured seed value, or NoSimulate if simulation is disabled.
// When the random condition is met, it delays the request by a duration defined in the configuration;
// otherwise, it passes the request to the next handler without delay.
func getSimulateSlowRequestsMiddleware(c *config.Config) func(http.Handler) http.Handler {
	if c.SimulateSlowRequests {
		return simulateSlowRequestsMiddleware(c)
	}
	return noSimulate
}

// SimulateErrorsMiddleware returns a middleware that randomly simulates server errors
// based on the configured seed value. When the random condition is met, it responds
// with a 500 Internal Server Error; otherwise, it passes the request to the next handler.

func simulateErrorsMiddleware(c *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c.SimulateErrors && rand.Intn(c.SeedForSimulateErrors) == 0 {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// getSimulateErrorsMiddleware returns a middleware that randomly simulates server errors
// based on the configured seed value, or NoSimulate if simulation is disabled.
// When the random condition is met, it responds with a 500 Internal Server Error;
// otherwise, it passes the request to the next handler.
func getSimulateErrorsMiddleware(c *config.Config) func(http.Handler) http.Handler {
	if c.SimulateErrors {
		return simulateErrorsMiddleware(c)
	}
	return noSimulate
}

// getBucketLimiterMiddleware returns a middleware that uses a TokenBucket to rate limit incoming requests.
// It allows a certain number of requests to pass through the rate limiter, and returns a 429 Too Many Requests
// status code if the request is not allowed. If the given configuration disables rate limiting, it returns
// a middleware that passes all requests to the next handler.
func getBucketLimiterMiddleware(c *config.Config) func(http.Handler) http.Handler {
	if c.QtTokens > 0 && c.TimeFrame > 0 {
		l := limiter.NewTokenBucket(c.QtTokens, c.TimeFrame)
		return l.Middleware
	}
	return noSimulate
}

// getAllMiddlewares returns a middleware that composes all the other middlewares in the order of rate limiting,
// simulating slow requests and simulating errors. It takes a configuration object as argument and returns a middleware
// that takes a next http.Handler argument and returns a new http.Handler that will be used to handle the next request.
func getAllMiddlewares(c *config.Config) func(http.Handler) http.Handler {
	b := getBucketLimiterMiddleware(c)
	s := getSimulateSlowRequestsMiddleware(c)
	e := getSimulateErrorsMiddleware(c)
	return func(next http.Handler) http.Handler {
		return b(s(e(next)))
	}
}
