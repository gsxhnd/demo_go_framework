package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type limiterEntry struct {
	limiter    *rate.Limiter
	lastSeenAt time.Time
}

type limiterStore struct {
	mu      sync.RWMutex
	entries map[string]*limiterEntry
	rps     rate.Limit
	burst   int
	ttl     time.Duration
}

func newLimiterStore(rps rate.Limit, burst int, ttl time.Duration) *limiterStore {
	return &limiterStore{
		entries: make(map[string]*limiterEntry),
		rps:     rps,
		burst:   burst,
		ttl:     ttl,
	}
}

func (s *limiterStore) get(key string) *rate.Limiter {
	now := time.Now()

	s.mu.RLock()
	entry, ok := s.entries[key]
	s.mu.RUnlock()
	if ok {
		s.mu.Lock()
		entry.lastSeenAt = now
		s.mu.Unlock()
		return entry.limiter
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok = s.entries[key]
	if ok {
		entry.lastSeenAt = now
		return entry.limiter
	}

	limiter := rate.NewLimiter(s.rps, s.burst)
	s.entries[key] = &limiterEntry{
		limiter:    limiter,
		lastSeenAt: now,
	}

	return limiter
}

func (s *limiterStore) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanup()
	}
}

func (s *limiterStore) cleanup() {
	deadline := time.Now().Add(-s.ttl)

	s.mu.Lock()
	defer s.mu.Unlock()

	for key, entry := range s.entries {
		if entry.lastSeenAt.Before(deadline) {
			delete(s.entries, key)
		}
	}
}

func keyFromRequest(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	return host + ":" + r.Method + ":" + r.URL.Path
}

func newRateLimitMiddleware(store *limiterStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFromRequest(r)
			limiter := store.get(key)

			if !limiter.Allow() {
				w.Header().Set("Retry-After", "1")
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func main() {
	store := newLimiterStore(rate.Limit(2), 4, 10*time.Minute)
	go store.cleanupLoop(time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, limited world\n")
	})

	handler := chain(mux, newRateLimitMiddleware(store))

	addr := ":8081"
	log.Printf("rate limit example listening on %s", addr)
	log.Printf("try: for i in {1..10}; do curl -i http://localhost%s/hello; done", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
