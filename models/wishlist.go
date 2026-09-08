package models

import "time"

type WishlistItem struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"userId" gorm:"index;not null"`
	ProductID string    `json:"productId" gorm:"index;not null"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID;references:ID"`
	CreatedAt time.Time `json:"createdAt"`
}

type ToggleWishlistRequest struct {
	ProductID string `json:"productId" binding:"required"`
}
