package middleware

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/stdlib"
	"github.com/ulule/limiter/drivers/store/memory"
)

// RateLimiter represents an in memory rate limiter
type RateLimiter struct {
	*stdlib.Middleware
}

const (
	// DefaultRateLimitPeriod is the default time window to track calls
	DefaultRateLimitPeriod = time.Minute * 10
	// DefaultRateLimit is the default amount of allowed calls to the protected api in the time window
	DefaultRateLimit = 50
)

// RateLimit creates a new rate limiting middleware with in memory store
// the amount of maximum requests can be set as wel as the duration in which this limit
// applies
func RateLimit(period time.Duration, limit int) RateLimiter {
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: period,
		Limit:  int64(limit),
	}

	lmt := limiter.New(store, rate)
	middleware := stdlib.NewMiddleware(lmt, stdlib.WithForwardHeader(true))
	middleware.OnLimitReached = func(w http.ResponseWriter, r *http.Request) {
		// Get the clients ip address.
		ipString := r.RemoteAddr
		// Account for proxies.
		if forwardString := r.Header.Get("X-Forwarded-For"); forwardString != "" {
			ipString = forwardString
		}
		// Cloudflare creates a special ipv6 address if the client uses ipv6,
		// so check for the appropriate header
		// We will rate limit on this special ipv4 address, but this is not a problem
		// as Cloudflare maps these special addresses to the underlying ipv6 address
		// in a one on one relationship
		if ipv6 := r.Header.Get("Cf-Connecting-Ipv6"); ipv6 != "" {
			ipString = ipv6
		}
		log.Info("Rate limiting request from: ", ipString)

		// Write some info back to the client
		w.WriteHeader(http.StatusTooManyRequests)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("You have reached the maximum request limit."))
		return
	}
	return RateLimiter{middleware}

}
