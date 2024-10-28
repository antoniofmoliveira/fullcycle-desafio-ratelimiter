package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	Port                        int
	QtTokens                    int
	TimeFrame                   time.Duration
	SimulateSlowRequests        bool
	SimulateErrors              bool
	SeedForSimulateErrors       int
	SeedForSimulateSlowRequests int
}

// NewConfig parses flags and creates a new Config object.
// It returns a Config object with the parsed flags.
// If the flags are invalid, it logs the error and exits the program.
func NewConfig() *Config {

	port := flag.Int("port", 8080, "Port to listen on")
	qtTokens := flag.Int("qt-tokens", 100000, "Number of tokens in the token bucket. 0 to disable rate limiting")
	timeFrameSeconds := flag.Int("time-frame-seconds", 1, "Time frame in seconds. 0 to disable rate limiting")
	simulateErrors := flag.Bool("simulate-errors", false, "Simulate errors")
	simulateSlowRequests := flag.Bool("simulate-slow-requests", false, "Simulate slow requests")
	seedForSimulateErrors := flag.Int("seed-for-simulate-errors", 100, "Seed for simulate errors. 100 equals 1 in 100 requests randomly fail")
	seedForSimulateSlowRequests := flag.Int("seed-for-simulate-slow-requests", 100, "Seed for simulate slow requests (microseconds)")
	flag.Parse()

	c := &Config{
		Port:                        *port,
		QtTokens:                    *qtTokens,
		TimeFrame:                   time.Duration(*timeFrameSeconds) * time.Second,
		SimulateSlowRequests:        *simulateSlowRequests,
		SimulateErrors:              *simulateErrors,
		SeedForSimulateErrors:       *seedForSimulateErrors,
		SeedForSimulateSlowRequests: *seedForSimulateSlowRequests,
	}

	if err := c.Validate(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	return c
}

var (
	ErrInvalidPort                        = fmt.Errorf("invalid port")
	ErrInvalidQtTokens                    = fmt.Errorf("invalid qt tokens")
	ErrInvalidTimeFrame                   = fmt.Errorf("invalid time frame")
	ErrInvalidSeedForSimulateErrors       = fmt.Errorf("invalid seed for simulate errors")
	ErrInvalidSeedForSimulateSlowRequests = fmt.Errorf("invalid seed for simulate slow requests")
)

// Validate checks the configuration values and returns an error if any of them are invalid.
// It validates the following conditions:
// - The port must be between 1 and 65535.
// - The number of tokens (QtTokens) must be non-negative.
// - The time frame must be non-negative.
// - The seed for simulating errors must be non-negative.
// - The seed for simulating slow requests must be non-negative.
// If any validation fails, it logs the error and returns the corresponding error.
func (c *Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		slog.Error(ErrInvalidPort.Error(), "port", c.Port)
		return ErrInvalidPort
	}
	if c.QtTokens < 0 {
		slog.Error(ErrInvalidQtTokens.Error(), "qtTokens", c.QtTokens)
		return ErrInvalidQtTokens
	}
	if c.TimeFrame < 0 {
		slog.Error(ErrInvalidTimeFrame.Error(), "timeFrame", c.TimeFrame)
		return ErrInvalidTimeFrame
	}
	if c.SeedForSimulateErrors < 0 {
		slog.Error(ErrInvalidSeedForSimulateErrors.Error(), "seedForSimulateErrors", c.SeedForSimulateErrors)
		return ErrInvalidSeedForSimulateErrors
	}
	if c.SeedForSimulateSlowRequests < 0 {
		slog.Error(ErrInvalidSeedForSimulateSlowRequests.Error(), "seedForSimulateSlowRequests", c.SeedForSimulateSlowRequests)
		return ErrInvalidSeedForSimulateSlowRequests
	}
	return nil
}
