package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"order-service/middleware"
	"order-service/models"
	"os"
	"shared/jwt"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func getDB() (*sql.DB, error) {
	connStr := os.Getenv("ORDER_DB_URL")
	return sql.Open("postgres", connStr)
}

func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*jwt.Claims)
	if !ok || claims == nil {
		log.Println("[PlaceOrder] Invalid or missing claims in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	db, err := getDB()
	if err != nil {
		http.Error(w, "DB connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	order.Date = time.Now()
	var orderID int
	err = db.QueryRow(`
		INSERT INTO orders (email, mobile, name, address, district, state, country, pincode, date, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'Placed') RETURNING id`,
		order.Email, order.Mobile, order.Name, order.Address, order.District, order.State, order.Country, order.Pincode, order.Date,
	).Scan(&orderID)
	if err != nil {
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	for _, item := range order.Items {
		_, err := db.Exec(`
			INSERT INTO order_items (order_id, name, price, qty, image, seller_id)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			orderID, item.Name, item.Price, item.Qty, item.Image, item.SellerID,
		)
		if err != nil {
			http.Error(w, "Item insert failed", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order placed successfully"})
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*jwt.Claims)
	if !ok || claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db, err := getDB()
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT o.id, o.date, o.status, i.name, i.price, i.qty, i.image
		FROM orders o
		JOIN order_items i ON o.id = i.order_id
		WHERE o.email = $1 OR o.mobile = $2
		ORDER BY o.date DESC`, claims.Email, claims.Mobile)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	orderMap := make(map[int]*models.Order)
	for rows.Next() {
		var oid int
		var date time.Time
		var status string
		var item models.CartItem
		err := rows.Scan(&oid, &date, &status, &item.Name, &item.Price, &item.Qty, &item.Image)
		if err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		if orderMap[oid] == nil {
			orderMap[oid] = &models.Order{ID: oid, Date: date, Status: status}
		}
		orderMap[oid].Items = append(orderMap[oid].Items, item)
	}

	var orders []models.Order
	for _, o := range orderMap {
		orders = append(orders, *o)
	}
	json.NewEncoder(w).Encode(orders)
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*jwt.Claims)
	if !ok || claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}

	db, err := getDB()
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	res, err := db.Exec(`
		UPDATE orders
		SET status = 'Cancelled'
		WHERE id = $1 AND (email = $2 OR mobile = $3)`, orderID, claims.Email, claims.Mobile)
	if err != nil {
		http.Error(w, "Cancel failed", http.StatusInternalServerError)
		return
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "Order not found or unauthorized", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Order cancelled successfully"})
}

func GetOrdersBySeller(w http.ResponseWriter, r *http.Request) {
	sellerID := mux.Vars(r)["seller_id"]

	db, err := getDB()
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `
		SELECT 
			o.id, o.date, o.status,
			o.email, o.name, o.mobile, o.address, o.district, o.state, o.country, o.pincode,
			i.name, i.price, i.qty, i.image, i.seller_id
		FROM orders o
		JOIN order_items i ON o.id = i.order_id
		WHERE i.seller_id = $1
		ORDER BY o.date DESC
	`

	rows, err := db.Query(query, sellerID)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	orderMap := make(map[int]*models.Order)
	for rows.Next() {
		var oid int
		var date time.Time
		var status string
		var email, name, mobile, address, district, state, country, pincode string
		var item models.CartItem

		err := rows.Scan(
			&oid, &date, &status,
			&email, &name, &mobile, &address, &district, &state, &country, &pincode,
			&item.Name, &item.Price, &item.Qty, &item.Image, &item.SellerID,
		)
		if err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}

		if orderMap[oid] == nil {
			orderMap[oid] = &models.Order{
				ID:      oid,
				Date:    date,
				Status:  status,
				Email:   email,
				Name:    name,
				Mobile:  mobile,
				Address: address,
				District: district,
				State:    state,
				Country:  country,
				Pincode:  pincode,
			}
		}

		orderMap[oid].Items = append(orderMap[oid].Items, item)
	}

	var orders []models.Order
	for _, o := range orderMap {
		orders = append(orders, *o)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

