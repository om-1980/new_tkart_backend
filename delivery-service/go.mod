module delivery-service

go 1.24.4

replace shared => ../shared

require github.com/labstack/echo/v4 v4.13.4

require golang.org/x/time v0.11.0 // indirect

require github.com/golang-jwt/jwt/v5 v5.2.3 // indirect

require (
	github.com/joho/godotenv v1.5.1
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/lib/pq v1.10.9
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	shared v0.0.0-00010101000000-000000000000
)
