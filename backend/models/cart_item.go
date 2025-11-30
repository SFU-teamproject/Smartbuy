package models

type CartItem struct {
	ID           int `json:"id"`
	CartID       int `json:"cart_id"`
	SmartphoneID int `json:"smartphone_id"`
	Quantity     int `json:"quantity"`
}

type CartItemRequest struct {
	SmartphoneID int `json:"smartphone_id"`
	Quantity     int `json:"quantity"`
}

type SetItemQuantity struct {
	Quantity int `json:"quantity"`
}
