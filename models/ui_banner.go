package models

import "time"

type UIBanner struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	ImageURL  string    `json:"imageUrl" gorm:"not null"`
	MediaType string    `json:"mediaType" gorm:"default:'image'"` // 'image' or 'video'
	ActionURL string    `json:"actionUrl"`
	Tag       string    `json:"tag"` // e.g. 'SUMMER SALE', 'HERO', 'PROMO'
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	SortOrder int       `json:"sortOrder" gorm:"default:0"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateUIBannerRequest struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	ImageURL  string `json:"imageUrl" binding:"required"`
	MediaType string `json:"mediaType"`
	ActionURL string `json:"actionUrl"`
	Tag       string `json:"tag"`
	IsActive  bool   `json:"isActive"`
	SortOrder int    `json:"sortOrder"`
}
