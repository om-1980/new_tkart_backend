module order-service

go 1.24.4

replace shared => ../shared

require (
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/rs/cors v1.11.1
	shared v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.3 // indirect
	github.com/gorilla/mux v1.8.1
)
