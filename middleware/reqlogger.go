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
	log.Printf("\n\033[1;30;102mRequest headers:\033[0;0;0m\n" +
		"\033[1;36mAuthorization:\033[0;0m\n" +
		"%s\n" +
		"\033[1;36mUser-Agent:\033[0;0m\n" +
		"%s\n" +
		"\033[1;30;103mNetwork details:\033[0;0;0m\n" +
		"\033[1;32mRemote Addr:\033[0;0m\n" +
		"%s\n" +
		"\033[1;32mLocal Endpoint:\033[0;0m\n" +
		"%s\n" +
		"\033[1;32mRequest duration:\033[0;0m\n" +
		"%s\n",
		request.Header.Get("Authorization"),
		request.Header.Get("User-Agent"),
		request.RemoteAddr,
		request.URL.String(),
		time.Since(startReq).String(),
	)
}

func NewLogger(handler http.Handler) *Logger {
	return &Logger{handler}
}

