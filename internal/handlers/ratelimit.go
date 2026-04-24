package handlers

import (
	"net/http"
	"strings"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimiter limits requests per client IP using token buckets.
// The bucket map grows unbounded; Cloud Run instance recycling provides implicit cleanup.
type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*rate.Limiter
	rps     rate.Limit
	burst   int
}

// NewRateLimiter creates a limiter with the given requests-per-minute rate and burst size.
func NewRateLimiter(rpm float64, burst int) *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*rate.Limiter),
		rps:     rate.Limit(rpm / 60),
		burst:   burst,
	}
}

func (rl *RateLimiter) get(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	l, ok := rl.buckets[ip]
	if !ok {
		l = rate.NewLimiter(rl.rps, rl.burst)
		rl.buckets[ip] = l
	}
	return l
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if i := strings.LastIndex(ip, ":"); i != -1 {
			ip = ip[:i]
		}
		if !rl.get(ip).Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
