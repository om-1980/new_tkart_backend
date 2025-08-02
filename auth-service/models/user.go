package models

type User struct {
	ID            int    `json:"id"`
	SellerID      string `json:"seller_id"`
	DeliveryID	  string `json:"delivery_id"`
	Name          string `json:"name"`
	Password      string `json:"password,omitempty"`
	Email         string `json:"email,omitempty"`
	Mobile        string `json:"mobile"`
	Address       string `json:"address,omitempty"`
	District      string `json:"district,omitempty"`
	State         string `json:"state,omitempty"`
	Country       string `json:"country,omitempty"`
	Pincode       string `json:"pincode,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	ProfilePhoto  string `json:"profile_photo,omitempty"`
	Role          string `json:"role"`
	IsActive	  bool	 `json:"is_active"`
}

type Credentials struct {
	Identifier string `json:"identifier"` // seller_id or email/mobile or delivery_id
	Password   string `json:"password"`
	Role       string `json:"role"`
}