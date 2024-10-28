package server

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/antoniofmoliveira/fullcycle-desafio-ratelimiter/internal/config"
	"github.com/antoniofmoliveira/fullcycle-desafio-ratelimiter/internal/limiter"
)

// NoSimulate is a middleware that simply forwards the request to the next handler
// in the chain, without performing any simulation. It is used when the
// simulation is disabled in the configuration.
func NoSimulate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

// SimulateSlowRequestsMiddleware returns a middleware that simulates slow requests
// based on the configured seed value. When the random condition is met, it delays
// the request by a duration defined in the configuration; otherwise, it passes
// the request to the next handler without delay.
func SimulateSlowRequestsMiddleware(c *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c.SimulateSlowRequests && rand.Intn(c.SeedForSimulateSlowRequests) == 0 {
				time.Sleep(time.Duration(rand.Intn(c.SeedForSimulateSlowRequests)) * time.Microsecond)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GetSimulateSlowRequestsMiddleware returns a middleware that simulates slow requests
// based on the configured seed value, or NoSimulate if simulation is disabled.
// When the random condition is met, it delays the request by a duration defined in the configuration;
// otherwise, it passes the request to the next handler without delay.
func GetSimulateSlowRequestsMiddleware(c *config.Config) func(http.Handler) http.Handler {
	if c.SimulateSlowRequests {
		return SimulateSlowRequestsMiddleware(c)
	}
	return NoSimulate
}

// SimulateErrorsMiddleware returns a middleware that randomly simulates server errors
// based on the configured seed value. When the random condition is met, it responds
// with a 500 Internal Server Error; otherwise, it passes the request to the next handler.

func SimulateErrorsMiddleware(c *config.Config) func(http.Handler) http.Handler {
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

// GetSimulateErrorsMiddleware returns a middleware that randomly simulates server errors
// based on the configured seed value, or NoSimulate if simulation is disabled.
// When the random condition is met, it responds with a 500 Internal Server Error;
// otherwise, it passes the request to the next handler.
func GetSimulateErrorsMiddleware(c *config.Config) func(http.Handler) http.Handler {
	if c.SimulateErrors {
		return SimulateErrorsMiddleware(c)
	}
	return NoSimulate
}

// GetBucketLimiterMiddleware returns a middleware that uses a TokenBucket to rate limit incoming requests.
// It allows a certain number of requests to pass through the rate limiter, and returns a 429 Too Many Requests
// status code if the request is not allowed. If the given configuration disables rate limiting, it returns
// a middleware that passes all requests to the next handler.
func GetBucketLimiterMiddleware(c *config.Config) func(http.Handler) http.Handler {
	if c.QtTokens > 0 && c.TimeFrame > 0 {
		l := limiter.NewTokenBucket(c.QtTokens, c.TimeFrame)
		return l.Middleware
	}
	return NoSimulate
}

// GetAllMiddlewares returns a middleware that composes all the other middlewares in the order of rate limiting,
// simulating slow requests and simulating errors. It takes a configuration object as argument and returns a middleware
// that takes a next http.Handler argument and returns a new http.Handler that will be used to handle the next request.
func GetAllMiddlewares(c *config.Config) func(http.Handler) http.Handler {
	b := GetBucketLimiterMiddleware(c)
	s := GetSimulateSlowRequestsMiddleware(c)
	e := GetSimulateErrorsMiddleware(c)
	return func(next http.Handler) http.Handler {
		return b(s(e(next)))
	}
}
