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

// ─────────────────────────────────────────────────────────────────────────────
// In-memory fallback (used when DB is unavailable)
// ─────────────────────────────────────────────────────────────────────────────

var memHeroVideo = models.HeroVideoConfig{
	ID:        1,
	VideoURL:  "", // Empty = customer app uses its bundled local asset
	IsEnabled: true,
	UpdatedAt: time.Now(),
}

// GetHeroVideo returns the current hero video configuration.
// GET /api/ui/hero-video   (public — used by customer app)
// GET /api/admin/hero-video (admin panel)
func GetHeroVideo(c *gin.Context) {
	if config.DB != nil {
		var cfg models.HeroVideoConfig
		result := config.DB.First(&cfg)
		if result.Error == nil {
			c.JSON(http.StatusOK, cfg)
			return
		}
	}
	c.JSON(http.StatusOK, memHeroVideo)
}

// UploadHeroVideo accepts a multipart video file, saves it under ./uploads/videos/,
// updates the hero_video_config record, and returns the new URL.
// POST /api/admin/hero-video/upload
func UploadHeroVideo(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No video file provided (field name: 'video')"})
		return
	}

	// Validate content type — accept common video MIME types
	ct := file.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "video/") && ct != "application/octet-stream" {
		// Fall back to extension check
		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := map[string]bool{".mp4": true, ".mov": true, ".webm": true, ".avi": true, ".mkv": true}
		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only video files are allowed (.mp4, .mov, .webm, .avi, .mkv)"})
			return
		}
	}

	// Ensure uploads/videos directory exists
	uploadDir := "./uploads/videos"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("hero_video_%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video file"})
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	videoURL := fmt.Sprintf("%s://%s/uploads/videos/%s", scheme, host, filename)

	// Update DB or in-memory
	cfg := models.HeroVideoConfig{
		ID:        1,
		VideoURL:  videoURL,
		IsEnabled: true,
		UpdatedAt: time.Now(),
	}

	if config.DB != nil {
		config.DB.Save(&cfg)
	}
	memHeroVideo = cfg

	c.JSON(http.StatusOK, gin.H{
		"message":  "Hero video uploaded successfully",
		"videoUrl": videoURL,
		"config":   cfg,
	})
}

// ToggleHeroVideo enables or disables the hero video without changing the URL.
// PUT /api/admin/hero-video/toggle
func ToggleHeroVideo(c *gin.Context) {
	if config.DB != nil {
		var cfg models.HeroVideoConfig
		config.DB.First(&cfg)
		cfg.IsEnabled = !cfg.IsEnabled
		cfg.UpdatedAt = time.Now()
		config.DB.Save(&cfg)
		c.JSON(http.StatusOK, gin.H{"message": "Hero video toggled", "config": cfg})
		return
	}
	memHeroVideo.IsEnabled = !memHeroVideo.IsEnabled
	memHeroVideo.UpdatedAt = time.Now()
	c.JSON(http.StatusOK, gin.H{"message": "Hero video toggled", "config": memHeroVideo})
}
