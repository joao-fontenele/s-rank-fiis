package model

import "time"

type Operation struct {
	ID           int
	Code         string
	Price        float64
	Amount       int
	PurchaseDate time.Time `json:"purchaseDate"`
}
