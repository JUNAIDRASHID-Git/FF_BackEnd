package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"funfillers/backend/config"
	"funfillers/backend/models"
	"github.com/gin-gonic/gin"
)

type UIController struct{}

func NewUIController() *UIController {
	return &UIController{}
}

var mockBanners = []models.UIBanner{
	{
		ID:        "banner_1",
		Title:     "Ultimate Gaming & Tech Sale ⚡️",
		Subtitle:  "Get up to 40% OFF on premium audio gear & mechanical keyboards",
		ImageURL:  "https://images.unsplash.com/photo-1542751371-adc38448a05e?w=1200",
		ActionURL: "/category/electronics",
		Tag:       "HERO BANNER",
		IsActive:  true,
		SortOrder: 1,
		CreatedAt: time.Now(),
	},
	{
		ID:        "banner_2",
		Title:     "Redefine Your Workstation Comfort",
		Subtitle:  "Ergonomic mesh chairs and smart height-adjustable standing desks",
		ImageURL:  "https://images.unsplash.com/photo-1524758631624-e2822e304c36?w=1200",
		ActionURL: "/category/furniture",
		Tag:       "FEATURED COLLECTION",
		IsActive:  true,
		SortOrder: 2,
		CreatedAt: time.Now(),
	},
}

func (uic *UIController) GetBanners(c *gin.Context) {
	if config.DB != nil {
		var banners []models.UIBanner
		if err := config.DB.Order("sort_order asc").Find(&banners).Error; err == nil && len(banners) > 0 {
			c.JSON(http.StatusOK, banners)
			return
		}
	}
	c.JSON(http.StatusOK, mockBanners)
}

func (uic *UIController) CreateBanner(c *gin.Context) {
	var req models.CreateUIBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	title := strings.TrimSpace(req.Title)
	if len(title) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Banner title must be at least 3 characters long"})
		return
	}

	imgURL := strings.TrimSpace(req.ImageURL)
	if imgURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Banner image URL is required"})
		return
	}

	banner := models.UIBanner{
		ID:        fmt.Sprintf("banner_%d", time.Now().UnixNano()/1e6),
		Title:     req.Title,
		Subtitle:  req.Subtitle,
		ImageURL:  req.ImageURL,
		ActionURL: req.ActionURL,
		Tag:       req.Tag,
		IsActive:  req.IsActive,
		SortOrder: req.SortOrder,
		CreatedAt: time.Now(),
	}

	if config.DB != nil {
		config.DB.Create(&banner)
	}
	mockBanners = append([]models.UIBanner{banner}, mockBanners...)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Banner created successfully",
		"banner":  banner,
	})
}

func (uic *UIController) ToggleBannerStatus(c *gin.Context) {
	id := c.Param("id")
	for i, b := range mockBanners {
		if b.ID == id {
			mockBanners[i].IsActive = !mockBanners[i].IsActive
			if config.DB != nil {
				config.DB.Model(&models.UIBanner{}).Where("id = ?", id).Update("is_active", mockBanners[i].IsActive)
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Banner status toggled",
				"banner":  mockBanners[i],
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Banner not found"})
}

func (uic *UIController) UpdateBanner(c *gin.Context) {
	id := c.Param("id")
	var req models.CreateUIBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update in DB
	if config.DB != nil {
		config.DB.Model(&models.UIBanner{}).Where("id = ?", id).Updates(map[string]interface{}{
			"title":      req.Title,
			"subtitle":   req.Subtitle,
			"image_url":  req.ImageURL,
			"action_url": req.ActionURL,
			"tag":        req.Tag,
			"is_active":  req.IsActive,
			"sort_order": req.SortOrder,
		})
	}

	// Update in-memory
	for i, b := range mockBanners {
		if b.ID == id {
			mockBanners[i].Title = req.Title
			mockBanners[i].Subtitle = req.Subtitle
			mockBanners[i].ImageURL = req.ImageURL
			mockBanners[i].ActionURL = req.ActionURL
			mockBanners[i].Tag = req.Tag
			mockBanners[i].IsActive = req.IsActive
			mockBanners[i].SortOrder = req.SortOrder
			c.JSON(http.StatusOK, gin.H{"message": "Banner updated", "banner": mockBanners[i]})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Banner not found"})
}

func (uic *UIController) DeleteBanner(c *gin.Context) {
	id := c.Param("id")
	if config.DB != nil {
		config.DB.Where("id = ?", id).Delete(&models.UIBanner{})
	}
	updated := make([]models.UIBanner, 0)
	for _, b := range mockBanners {
		if b.ID != id {
			updated = append(updated, b)
		}
	}
	mockBanners = updated
	c.JSON(http.StatusOK, gin.H{"message": "Banner deleted successfully"})
}

// UploadBannerImage accepts a multipart image file, saves it in ./uploads/banners/, and returns its public URL.
// POST /api/admin/ui/banners/upload
func (uic *UIController) UploadBannerImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		file, err = c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided (field name: 'image' or 'file')"})
			return
		}
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true, ".svg": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only image files are allowed (.jpg, .jpeg, .png, .webp, .gif, .svg)"})
		return
	}

	uploadDir := "./uploads/banners"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	filename := fmt.Sprintf("banner_%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image file"})
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	imageURL := fmt.Sprintf("%s://%s/uploads/banners/%s", scheme, host, filename)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Banner image uploaded successfully",
		"imageUrl": imageURL,
	})
}

