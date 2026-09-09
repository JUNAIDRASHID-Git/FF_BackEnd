package models

import "time"

type Address struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	UserID      string    `json:"userId" gorm:"index"`
	Title       string    `json:"title"`
	Label       string    `json:"label"`
	FullAddress string    `json:"fullAddress"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	Country     string    `json:"country"`
	Pincode     string    `json:"pincode"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	IsDefault   bool      `json:"isDefault" gorm:"default:true"`
	CreatedAt   time.Time `json:"createdAt"`
}
