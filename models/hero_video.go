package models

import "time"

// HeroVideoConfig stores the single hero video configuration for the storefront
type HeroVideoConfig struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	VideoURL  string    `json:"videoUrl" gorm:"not null"`
	IsEnabled bool      `json:"isEnabled" gorm:"default:true"`
	UpdatedAt time.Time `json:"updatedAt"`
}
