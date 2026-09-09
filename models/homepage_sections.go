package models

import "time"

type AnnouncementConfig struct {
	IsEnabled bool   `json:"isEnabled"`
	Text      string `json:"text"`
	BgColor   string `json:"bgColor"`
	TextColor string `json:"textColor"`
	ActionURL string `json:"actionUrl"`
}

type FlashSaleConfig struct {
	IsEnabled    bool      `json:"isEnabled"`
	Title        string    `json:"title"`
	Subtitle     string    `json:"subtitle"`
	DiscountText string    `json:"discountText"`
	EndTime      time.Time `json:"endTime"`
	ActionURL    string    `json:"actionUrl"`
}

type FeaturedCategoriesConfig struct {
	IsEnabled   bool     `json:"isEnabled"`
	CategoryIDs []string `json:"categoryIds"`
}

type HomepageSectionsConfig struct {
	ID                 int                      `json:"id" gorm:"primaryKey"`
	Announcement       AnnouncementConfig       `json:"announcement" gorm:"embedded;embeddedPrefix:announcement_"`
	FlashSale          FlashSaleConfig          `json:"flashSale" gorm:"embedded;embeddedPrefix:flash_sale_"`
	FeaturedCategories FeaturedCategoriesConfig `json:"featuredCategories" gorm:"serializer:json"`
	UpdatedAt          time.Time                `json:"updatedAt"`
}
