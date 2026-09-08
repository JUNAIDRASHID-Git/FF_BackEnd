package models

import "time"

type Category struct {
	ID            string        `json:"id" gorm:"primaryKey"`
	Name          string        `json:"name" gorm:"not null"`
	Slug          string        `json:"slug"`
	Icon          string        `json:"icon"`
	Image         string        `json:"image"`
	Description   string        `json:"description"`
	ItemCount     int           `json:"itemCount" gorm:"-"`
	SubCategories []SubCategory `json:"subCategories" gorm:"foreignKey:CategoryID"`
	CreatedAt     time.Time     `json:"createdAt"`
}

type SubCategory struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	CategoryID  string    `json:"categoryId" gorm:"index"`
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug"`
	Icon        string    `json:"icon"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateCategoryRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Icon        string `json:"icon"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

type CreateSubCategoryRequest struct {
	ID          string `json:"id"`
	CategoryID  string `json:"categoryId" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Icon        string `json:"icon"`
	Image       string `json:"image"`
	Description string `json:"description"`
}
