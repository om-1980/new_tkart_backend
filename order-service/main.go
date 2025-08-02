package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"order-service/handlers"
	"order-service/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	r := mux.NewRouter()
	r.Handle("/orders/place", middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.PlaceOrder))).Methods("POST")
	r.Handle("/orders", middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.GetOrders))).Methods("GET")
	r.Handle("/orders/cancel", middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.CancelOrder))).Methods("POST")
	r.Handle("/seller/orders/{seller_id}", middleware.SellerAuthMiddleware(http.HandlerFunc(handlers.GetOrdersBySeller))).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGIN")},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(r)

	log.Printf("Order service running on port %s", port)
	http.ListenAndServe(":"+port, corsHandler)
}
