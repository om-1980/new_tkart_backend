package models

type Product struct {
	ID          		int     `json:"id"`
	SellerID    		string  `json:"seller_id"`
	Name        		string  `json:"name"`
	Category			string	`json:"category"`
	Subcategory			string	`json:"subcategory"`
	InnerSubcategory	string	`json:"inner_subcategory"`
	Description 		string  `json:"description,omitempty"`
	Price       		float64 `json:"price"`
	Quantity    		int     `json:"quantity"`
	InStock				bool	`json:"in_stock"`
	Image1    			string  `json:"image1,omitempty"`
	Image2    			string  `json:"image2,omitempty"`
	Image3    			string  `json:"image3,omitempty"`
	Image4    			string  `json:"image4,omitempty"`
}