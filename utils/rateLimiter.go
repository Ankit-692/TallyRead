package utils

import (
	"sync"

	"golang.org/x/time/rate"
)

var (
	limiters = make(map[int64]*rate.Limiter)
	mu       sync.Mutex
)

func GetLimiter(userID int64) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	l, exists := limiters[userID]
	if !exists {
		// Allow 5 requests per second with a "burst" of 10
		l = rate.NewLimiter(5, 10)
		limiters[userID] = l
	}
	return l
}
