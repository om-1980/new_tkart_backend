package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"product-service/models"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func getDB() (*sql.DB, error) {
	connStr := os.Getenv("PRODUCT_DB_URL")
	return sql.Open("postgres", connStr)
}

func nullify(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	db, err := getDB()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid product data", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO products (seller_id, name, category, subcategory, inner_subcategory, description, price, quantity, in_stock, image1, image2, image3, image4) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err = db.Exec(query, product.SellerID, product.Name, product.Category, product.Subcategory, product.InnerSubcategory, product.Description, product.Price, product.Quantity, product.InStock, nullify(product.Image1), nullify(product.Image2), nullify(product.Image3), nullify(product.Image4))
	if err != nil {
		log.Println("DB Exec Error:", err)
		http.Error(w, "Failed to add product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product added successfully"})
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db, err := getDB()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, seller_id, name, category, subcategory, inner_subcategory, description, price, quantity, in_stock, image1, image2, image3, image4 FROM products`)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var img1, img2, img3, img4 []byte
		var p models.Product
		err := rows.Scan(&p.ID, &p.SellerID, &p.Name, &p.Category, &p.Subcategory, &p.InnerSubcategory, &p.Description, &p.Price, &p.Quantity, &p.InStock, &img1, &img2, &img3, &img4)
		if err != nil {
			http.Error(w, "Error reading product", http.StatusInternalServerError)
			return
		}
		if len(img1) > 0 {
			p.Image1 = string(img1)
		}
		if len(img2) > 0 {
			p.Image2 = string(img2)
		}
		if len(img3) > 0 {
			p.Image3 = string(img3)
		}
		if len(img4) > 0 {
			p.Image4 = string(img4)
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func GetProductsBySellerID(w http.ResponseWriter, r *http.Request) {
	db, err := getDB()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	sellerID := mux.Vars(r)["seller_id"]
	rows, err := db.Query(`
		SELECT id, seller_id, name, category, subcategory, inner_subcategory, description, price, quantity, in_stock, image1, image2, image3, image4 
		FROM products WHERE seller_id = $1`, sellerID)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var img1, img2, img3, img4 []byte

		err := rows.Scan(&p.ID, &p.SellerID, &p.Name, &p.Category, &p.Subcategory, &p.InnerSubcategory, &p.Description,
			&p.Price, &p.Quantity, &p.InStock, &img1, &img2, &img3, &img4)
		if err != nil {
			http.Error(w, "Error scanning product", http.StatusInternalServerError)
			return
		}

		if len(img1) > 0 {
			p.Image1 = string(img1)
		}
		if len(img2) > 0 {
			p.Image2 = string(img2)
		}
		if len(img3) > 0 {
			p.Image3 = string(img3)
		}
		if len(img4) > 0 {
			p.Image4 = string(img4)
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error reading rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}


func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	db, err := getDB()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]
	seller_id := vars["seller_id"]
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid product data", http.StatusBadRequest)
		return
	}

	query := `UPDATE products SET name=$1, category=$2, subcategory=$3, inner_subcategory=$4,description=$5, price=$6, quantity=$7, in_stock=$8, image1=$9, image2=$10, image3=$11, image4=$12 WHERE id=$13 AND seller_id=$14`
	result, err := db.Exec(query, product.Name, product.Category, product.Subcategory, product.InnerSubcategory, product.Description, product.Price, product.Quantity, product.InStock, product.Image1, product.Image2, product.Image3, product.Image4, id, seller_id)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking update result", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "No product found for given ID and seller", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Product updated"})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	db, err := getDB()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]
	seller_id := vars["seller_id"]
	_, err = db.Exec(`DELETE FROM products WHERE id=$1 AND seller_id=$2`, id, seller_id)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted"})
}