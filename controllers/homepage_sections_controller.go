package controllers

import (
	"net/http"
	"time"

	"funfillers/backend/config"
	"funfillers/backend/models"
	"github.com/gin-gonic/gin"
)

var memHomepageSections = models.HomepageSectionsConfig{
	ID: 1,
	Announcement: models.AnnouncementConfig{
		IsEnabled: true,
		Text:      "🎉 Free Express Delivery on orders over $49! Use code FUN50",
		BgColor:   "#6C5CE7",
		TextColor: "#FFFFFF",
		ActionURL: "/category/toys",
	},
	FlashSale: models.FlashSaleConfig{
		IsEnabled:    true,
		Title:        "⚡️ 24-HOUR FLASH SALE",
		Subtitle:     "Flat 50% OFF on Top Rated Action Figures & Plushies",
		DiscountText: "UP TO 50% OFF",
		EndTime:      time.Now().Add(12 * time.Hour),
		ActionURL:    "/category/toys",
	},
	FeaturedCategories: models.FeaturedCategoriesConfig{
		IsEnabled:   true,
		CategoryIDs: []string{"cat_1", "cat_2", "cat_3"},
	},
	UpdatedAt: time.Now(),
}

func GetHomepageSections(c *gin.Context) {
	if config.DB != nil {
		var cfg models.HomepageSectionsConfig
		if err := config.DB.First(&cfg).Error; err == nil {
			c.JSON(http.StatusOK, cfg)
			return
		}
	}
	c.JSON(http.StatusOK, memHomepageSections)
}

func UpdateHomepageSections(c *gin.Context) {
	var req models.HomepageSectionsConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = 1
	req.UpdatedAt = time.Now()

	if config.DB != nil {
		config.DB.Save(&req)
	}
	memHomepageSections = req

	c.JSON(http.StatusOK, gin.H{
		"message": "Homepage sections updated successfully",
		"config":  memHomepageSections,
	})
}
