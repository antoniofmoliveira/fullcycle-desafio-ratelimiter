package limiter

import (
	"log/slog"
	"net/http"
	"time"
)

// qtTokens: The total number of tokens in the token bucket.
// qtTokensRemain: The remaining number of tokens in the token bucket.
// timeFrame: The duration of the time frame for the token bucket.
// lastLimiterStart: The timestamp of the last time the token bucket was reset.
type tokenBucket struct {
	qtTokens         int
	qtTokensRemain   int
	timeFrame        time.Duration
	lastLimiterStart time.Time
}

// NewTokenBucket returns a new TokenBucket with the given qtTokens and timeFrame.
// A goroutine is started that will reset the token bucket every timeFrame duration.
// The bucket is reset immediately when NewTokenBucket is called, and then every timeFrame duration after that.
func NewTokenBucket(qtTokens int, timeFrame time.Duration) *tokenBucket {
	tb := &tokenBucket{
		qtTokens:         qtTokens,
		qtTokensRemain:   qtTokens,
		timeFrame:        timeFrame,
		lastLimiterStart: time.Now(),
	}

	go func() {
		for {
			time.Sleep(timeFrame)
			slog.Info("Reset bucket", "Tokens remaining", tb.qtTokensRemain, "Reset bucket at", time.Now())
			tb.ResetBucket()
		}
	}()

	return tb
}

// Allow returns true if a request is allowed to pass through the rate limiter,
// and false otherwise. It decrements the qtTokensRemain counter and returns
// true if it is greater than 0, and false otherwise.
func (t *tokenBucket) Allow() bool {
	if t.qtTokensRemain > 0 {
		t.qtTokensRemain--
		return true
	}
	return false
}

// ResetBucket resets the token bucket by setting the remaining tokens to the total tokens
// and updating the timestamp of the last bucket reset to the current time.
func (t *tokenBucket) ResetBucket() {
	t.qtTokensRemain = t.qtTokens
	t.lastLimiterStart = time.Now()
}

// TimeFrame returns the duration of the time frame for the token bucket.
func (t *tokenBucket) TimeFrame() time.Duration {
	return t.timeFrame
}

// LastLimiterStart returns the timestamp of the last time the token bucket was reset.
func (t *tokenBucket) LastLimiterStart() time.Time {
	return t.lastLimiterStart
}

// QtTokens returns the total number of tokens in the token bucket.
func (t *tokenBucket) QtTokens() int {
	return t.qtTokens
}

// Middleware returns a middleware that uses the TokenBucket to rate limit incoming requests.
// It writes a http.StatusTooManyRequests (429) status code if the request is not allowed
// to pass through the rate limiter, and calls the next handler in the chain
// otherwise.
func (t *tokenBucket) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !t.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
