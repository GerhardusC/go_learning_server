package middleware

import (
	"log"
	"net/http"
	"time"
)

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	startReq := time.Now()
	l.handler.ServeHTTP(writer, request)
	log.Printf(
		`Request headers:
		%#v",
		Endpoint:
		%s,
		Request duration:
		%v`,
		request.Header,
		request.URL,
		time.Since(startReq),
	)
}

func NewLogger(handler http.Handler) *Logger {
	return &Logger{handler}
}

