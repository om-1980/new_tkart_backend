package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"auth-service/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/info", handlers.GetUserInfo)
	mux.HandleFunc("/update", handlers.UpdateUser)
	mux.HandleFunc("/delete", handlers.DeleteUser)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("BUYER_ALLOWED_ORIGIN"), os.Getenv("SELLER_ALLOWED_ORIGIN")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(mux)

	log.Printf("Auth service is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
