package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"auth-service/models"
	"auth-service/utils"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

func getDBByRole(role string) (*sql.DB, error) {
	var connStr string
	switch role {
		case "buyer":
			connStr = os.Getenv("BUYER_DATABASE_URL")
		case "seller":
			connStr = os.Getenv("SELLER_DATABASE_URL")
		case "delivery":
			connStr = os.Getenv("DELIVERMAN_DATABASE_URL")
		default:
			return nil, fmt.Errorf("invalid role")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

func nullify(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Password == "" || user.Role == "" ||
		(user.Role == "buyer" && user.Email == "" && user.Mobile == "") ||
		(user.Role == "seller" && user.SellerID == "") ||
		(user.Role == "delivery" && user.DeliveryID == "") {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	db, err := getDBByRole(user.Role)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	switch user.Role {
		case "seller":
			query := `INSERT INTO sellers (seller_id, name, password, email, mobile, account_number, address, district, state, country, pincode, profile_photo, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, TRUE)`
			_, err = db.Exec(query, user.SellerID, user.Name, string(hashedPassword), user.Email, user.Mobile, user.AccountNumber, user.Address, user.District, user.State, user.Country, user.Pincode, user.ProfilePhoto)
		case "buyer":
			query := `INSERT INTO buyers (name, password, email, mobile, address, district, state, country, pincode, profile_photo) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
			_, err = db.Exec(query, user.Name, string(hashedPassword), user.Email, user.Mobile, user.Address, user.District, user.State, user.Country, user.Pincode, user.ProfilePhoto)
		case "delivery":
			query := `INSERT INTO deliveryman (delivery_id, name, password, email, mobile, account_number, address, district, state, country, pincode, profile_photo, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, TRUE)`
			_, err = db.Exec(query, user.DeliveryID, user.Name, string(hashedPassword), user.Email, user.Mobile, user.AccountNumber, user.Address, user.District, user.State, user.Country, user.Pincode, user.ProfilePhoto)
	}

	if err != nil {
		http.Error(w, "User registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	db, err := getDBByRole(creds.Role)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user models.User

	switch creds.Role {
		case "seller":
			query := `SELECT id, seller_id, name, password, email, mobile, is_active FROM sellers WHERE seller_id=$1`
			err = db.QueryRow(query, creds.Identifier).Scan(&user.ID, &user.SellerID, &user.Name, &user.Password, &user.Email, &user.Mobile, &user.IsActive)
			if err != nil {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}
			if !user.IsActive {
				http.Error(w, "Seller account is inactive", http.StatusForbidden)
				return
			}
		case "buyer":
			query := `SELECT id, name, password, email, mobile FROM buyers WHERE email=$1 OR mobile=$1`
			err = db.QueryRow(query, creds.Identifier).Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Mobile)
		case "delivery":
			query := `SELECT id, delivery_id, name, password, email, mobile, is_active FROM deliveryman WHERE delivery_id=$1`
			err = db.QueryRow(query, creds.Identifier).Scan(&user.ID, &user.DeliveryID, &user.Name, &user.Password, &user.Email, &user.Mobile, &user.IsActive)
			if err != nil {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}
			if !user.IsActive {
				http.Error(w, "Delivery account is inactive", http.StatusForbidden)
				return
			}
	}

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	var token string
	switch creds.Role {
		case "seller":
			token, err = utils.GenerateJWT(user.ID, "seller", user.SellerID, user.Name, user.Email, user.Mobile)
		case "buyer":
			identifier := user.Email
			if identifier == "" {
				identifier = user.Mobile
			}
			token, err = utils.GenerateJWT(user.ID, "buyer", identifier, user.Name, user.Email, user.Mobile)
		case "delivery":
			token, err = utils.GenerateJWT(user.ID, "delivery", user.DeliveryID, user.Name, user.Email, user.Mobile)
	}

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	user.Password = ""

	response := map[string]interface{}{
		"token": token,
		"name": user.Name,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	role := r.URL.Query().Get("role")
	identifier := r.URL.Query().Get("identifier") // seller_id, email, or mobile

	if role == "" || identifier == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	db, err := getDBByRole(role)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user models.User

	if role == "seller" {
		query := `
			SELECT id, seller_id, name, email, mobile, account_number, address, district, state, country, pincode, profile_photo, is_active
			FROM sellers
			WHERE seller_id = $1
		`
		err = db.QueryRow(query, identifier).Scan(
			&user.ID, &user.SellerID, &user.Name, &user.Email, &user.Mobile, &user.AccountNumber,
			&user.Address, &user.District, &user.State, &user.Country, &user.Pincode, &user.ProfilePhoto, &user.IsActive)
	} else if role == "buyer" {
		query := `
			SELECT id, name, email, mobile, address, district, state, country, pincode, profile_photo
			FROM buyers
			WHERE email = $1 OR mobile = $1
		`
		err = db.QueryRow(query, identifier).Scan(
			&user.ID, &user.Name, &user.Email, &user.Mobile,
			&user.Address, &user.District, &user.State, &user.Country, &user.Pincode, &user.ProfilePhoto)
	} else if role == "delivery" {
		query := `
			SELECT id, delivery_id, name, email, mobile, account_number, address, district, state, country, pincode, profile_photo, is_active
			FROM deliveryman
			WHERE delivery_id = $1
		`
		err = db.QueryRow(query, identifier).Scan(
			&user.ID, &user.DeliveryID, &user.Name, &user.Email, &user.Mobile, &user.AccountNumber,
			&user.Address, &user.District, &user.State, &user.Country, &user.Pincode, &user.ProfilePhoto, &user.IsActive)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			fmt.Println("Error fetching user:", err) // <== log actual error to server
			http.Error(w, "Error fetching user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}


func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	db, err := getDBByRole(user.Role)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	switch user.Role {
		case "seller":
			query := `UPDATE sellers SET name=$1, email=$2, account_number=$3, address=$4, district=$5, state=$6, country=$7, pincode=$8, profile_photo=$9, is_active=$10 WHERE id=$11 or seller_id=$12`
			_, err = db.Exec(query, user.Name, nullify(user.Email), user.AccountNumber, nullify(user.Address), nullify(user.District), nullify(user.State), nullify(user.Country), nullify(user.Pincode), user.ProfilePhoto, user.IsActive, user.ID, user.SellerID)
		case "buyer":
			query := `UPDATE buyers SET name=$1, email=$2, address=$3, district=$4, state=$5, country=$6, pincode=$7, profile_photo=$8 WHERE id=$9`
			_, err = db.Exec(query, user.Name, nullify(user.Email), nullify(user.Address), nullify(user.District), nullify(user.State), nullify(user.Country), nullify(user.Pincode), user.ProfilePhoto, user.ID)
		case "delivery":
			query := `UPDATE deliveryman SET name=$1, email=$2, account_number=$3, address=$4, district=$5, state=$6, country=$7, pincode=$8, profile_photo=$9, is_active=$10 WHERE id=$11 or delivery_id=$12`
			_, err = db.Exec(query, user.Name, nullify(user.Email), user.AccountNumber, nullify(user.Address), nullify(user.District), nullify(user.State), nullify(user.Country), nullify(user.Pincode), user.ProfilePhoto, user.IsActive, user.ID, user.DeliveryID)
	}

	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID   int    `json:"id"`
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	db, err := getDBByRole(req.Role)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var query string
	switch req.Role {
		case "seller":
			query = `DELETE FROM sellers WHERE id=$1`
		case "buyer":
			query = `DELETE FROM buyers WHERE id=$1`
		case "delivery":
			query = `DELETE FROM deliveryman WHERE id=$1`
	}

	_, err = db.Exec(query, req.ID)
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}