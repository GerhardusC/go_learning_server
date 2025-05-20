package middleware

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"net/http"

	"sync"
	"time"

	"golang.org/x/time/rate"
)

const BURST_REQUEST_COUNT int = 10

type rateLimitObj struct {
	limiter *rate.Limiter
	lastSeen time.Time
}

type clientsObj struct {
	mutx sync.Mutex
	clients map[string]*rateLimitObj
	channel chan string
}

func clearOldClients (clientsMapObj *clientsObj) {
	for {
		time.Sleep(time.Minute)

		for key, client := range clientsMapObj.clients {
			if time.Since(client.lastSeen) > time.Hour * 1 {
				clientsMapObj.mutx.Lock()
				delete(clientsMapObj.clients, key)
				clientsMapObj.mutx.Unlock()
				log.Println("\n\033[1;34;47m Client deleted: \033[0m", key)
			}
		}

	}
}

func handleIPsSentToChan(clientsMapObj *clientsObj, rateLimit float64, burst int) {
	for  {
		addr := <- clientsMapObj.channel
		newVal := rateLimitObj{
			limiter: rate.NewLimiter(rate.Limit(rateLimit), burst),
			lastSeen: time.Now(),
		}
		clientsMapObj.clients[addr] = &newVal
	}
}

func LimitRate (next http.HandlerFunc, rateLimit float64, burst int) http.HandlerFunc {

	clientsMapObject := clientsObj{
		clients: make(map[string]*rateLimitObj),
		channel: make(chan string, 10),
	}
	go handleIPsSentToChan(&clientsMapObject, 0.2, burst)
	go clearOldClients(&clientsMapObject)

	return func (writer http.ResponseWriter, request *http.Request) {
		addr := fmt.Sprintf("%x", sha256.Sum256([]byte(request.RemoteAddr)))

		val, ok := clientsMapObject.clients[addr]

		if !ok {
			clientsMapObject.channel<- addr
			next(writer, request)
			return
		}

		val.lastSeen = time.Now()

		if !val.limiter.Allow() {
			http.Error(
				writer,
				errors.New("Rate limit exceeded").Error(),
				http.StatusTooManyRequests,
			)
			return
		}

		next(writer, request)
	}
}

