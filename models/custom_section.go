package models

import "time"

type CustomSection struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	Title      string    `json:"title" gorm:"not null"`
	ProductIDs []string  `json:"productIds" gorm:"serializer:json"`
	IsActive   bool      `json:"isActive" gorm:"default:true"`
	SortOrder  int       `json:"sortOrder" gorm:"default:0"`
	CreatedAt  time.Time `json:"createdAt"`
}

type CreateCustomSectionRequest struct {
	Title      string   `json:"title" binding:"required"`
	ProductIDs []string `json:"productIds"`
	IsActive   bool     `json:"isActive"`
	SortOrder  int      `json:"sortOrder"`
}
