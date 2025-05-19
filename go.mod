module testing-server

go 1.24.2

require (
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/mattn/go-sqlite3 v1.14.28
	github.com/redis/go-redis/v9 v9.8.0
	golang.org/x/time v0.11.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

require (
	github.com/eclipse/paho.mqtt.golang v1.5.0 // direct
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
)
