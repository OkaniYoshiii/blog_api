package middleware

import (
	"github.com/didip/tollbooth/v8"
	"github.com/didip/tollbooth/v8/limiter"
)

// Creates a new rate limiter middleware that
// will check multiple headers to get IP
// address from the request
func RateLimiterMiddleware(rate float64) Middleware {
	const LookupsCount = 3

	lookups := [LookupsCount]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}
	middlewares := [LookupsCount]Middleware{}
	for i, name := range lookups {
		lmt := tollbooth.NewLimiter(rate, &limiter.ExpirableOptions{})
		lmt.SetIPLookup(limiter.IPLookup{
			Name:           name,
			IndexFromRight: 0,
		})

		middlewares[i] = tollbooth.HTTPMiddleware(lmt)
	}

	return Pipe(middlewares[:]...)
}
