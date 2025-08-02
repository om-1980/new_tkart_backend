package main

import (
	"log"
	"net/http"
	"os"

	"product-service/handlers"
	"product-service/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")

	r := mux.NewRouter()
	r.Use(middleware.JWTAuthMiddleware)
	r.HandleFunc("/products", handlers.AddProduct).Methods("POST")
	r.HandleFunc("/products", handlers.GetAllProducts).Methods("GET")
	r.HandleFunc("/products/{seller_id}", handlers.GetProductsBySellerID).Methods("GET")
	r.HandleFunc("/products/{id}/{seller_id}", handlers.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}/{seller_id}", handlers.DeleteProduct).Methods("DELETE")

	log.Printf("Product Service running on port %s", port)
	http.ListenAndServe(":"+port, r)
}
