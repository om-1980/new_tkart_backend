package models

type Delivery struct {
    DeliveryID      int    `json:"delivery_id"`
    OrderID         int    `json:"order_id"`
    DeliveryPerson  string `json:"delivery_person"`
    Status          string `json:"status"` // pending, picked, delivered, return_initiated, returned
    IsCOD           bool   `json:"is_cod"`
    CODAmount       int    `json:"cod_amount"`
    IsReturn        bool   `json:"is_return"`
}
