package middleware

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type rateLimitObj struct {
	limiter *rate.Limiter
	lastSeen time.Time
}


func LimitRate (next http.HandlerFunc, rateLimit float64) http.HandlerFunc {

	var (
		mutx sync.Mutex
		clients = make(map[string] *rateLimitObj)
	)

	return func (writer http.ResponseWriter, request *http.Request) {
		addr := fmt.Sprintf("%x", sha256.Sum256([]byte(request.RemoteAddr)))

		mutx.Lock()

		// TODO: Periodically clear map
		val, ok := clients[addr]

		if !ok {
			newVal := rateLimitObj{
				limiter: rate.NewLimiter(rate.Limit(rateLimit), 1),
			}
			clients[addr] = &newVal

			val = clients[addr]
		}

		val.lastSeen = time.Now()

		if !val.limiter.Allow() {
			http.Error(
				writer,
				errors.New("Rate limit exceeded").Error(),
				http.StatusTooManyRequests,
			)
			mutx.Unlock()
			return
		}
		mutx.Unlock()

		next(writer, request)
	}
}

