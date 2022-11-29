package prometheus_tools

import (
	"time"
)

func GenDuration(unit string, duration int) time.Duration {
	timeout := time.Millisecond * time.Duration(duration)
	if unit == "s" {
		timeout = time.Second * time.Duration(duration)
	} else if unit == "min" {
		timeout = time.Minute * time.Duration(duration)
	} else if unit == "ms" {
		timeout = time.Millisecond * time.Duration(duration)
	} else if unit == "mu" {
		timeout = time.Microsecond * time.Duration(duration)
	} else if unit == "ns" {
		timeout = time.Nanosecond * time.Duration(duration)
	}
	return timeout
}

// NewGlobalTicker create a time tick with specified duration
func NewGlobalTicker(unit string, duration int) *time.Ticker {
	timeout := GenDuration(unit, duration)
	timer := time.NewTicker(timeout)
	return timer
}
