package models

import "time"

type CartItem struct {
	Name  		string  `json:"name"`
	Price 		float64 `json:"price"`
	Qty  		int     `json:"qty"`
	Image 		string  `json:"image"`
	SellerID 	string 	`json:"seller_id"`
}

type Order struct {
	ID       int        `json:"id,omitempty"`
	Email    string     `json:"email,omitempty"`
	Mobile   string     `json:"mobile"`
	Name	 string 	`json:"name"`
	Address  string     `json:"address,omitempty"`
	District string		`json:"district,omitempty"`
	State    string 	`json:"state,omitempty"`
	Country  string 	`json:"country,omitempty"`
	Pincode  string     `json:"pincode"`
	Date     time.Time  `json:"date,omitempty"`
	Items    []CartItem `json:"items"`
	Status   string     `json:"status,omitempty"`
}
