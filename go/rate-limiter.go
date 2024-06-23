package main

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter controls how frequently events are allowed to happen
type RateLimiter struct {
	rate      int
	burst     int
	tokens    int
	lastCheck time.Time
	mu        sync.Mutex
}

// NewRateLimiter returns a new RateLimiter.
func NewRateLimter(rate int, burst int) *RateLimiter {
	return &RateLimiter{
		rate:      rate,
		burst:     burst,
		tokens:    burst,
		lastCheck: time.Now(),
	}
}

// Allow checks if a request is allowed and updates the token bucket
func (r1 *RateLimiter) Allow() bool {
	r1.mu.Lock()
	defer r1.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r1.lastCheck).Seconds()
	r1.lastCheck = now

	// Add new tokens based on the elapsed time
	r1.tokens += int(elapsed * float64(r1.rate))
	if r1.tokens > r1.burst {
		r1.tokens--
		return true
	}

	// Check if there are enough tokens for the request
	if r1.tokens > 0 {
		r1.tokens--
		return true
	}

	return false
}

func main() {
	rateLimter := NewRateLimter(5, 10) // 5 tokens per second, burst of 10 tokens

	for i := range 20 {
		if rateLimter.Allow() {
			fmt.Printf("Request %d: Allowed\n", i+1)
		} else {
			fmt.Printf("Request %d: Denied\n", i+1)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
