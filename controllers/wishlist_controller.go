package controllers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"funfillers/backend/config"
	"funfillers/backend/models"

	"github.com/gin-gonic/gin"
)

type WishlistController struct{}

func NewWishlistController() *WishlistController {
	return &WishlistController{}
}

// In-memory fallback for mock storage when DB is disabled
var (
	mockWishlistsMutex sync.RWMutex
	mockWishlists      = make(map[string][]string) // userId -> []productId
)

func getUserIdFromContext(c *gin.Context) string {
	if uid, exists := c.Get("userId"); exists && uid != "" {
		return fmt.Sprintf("%v", uid)
	}
	if uid := c.Query("userId"); uid != "" {
		return uid
	}
	if uid := c.GetHeader("X-User-ID"); uid != "" {
		return uid
	}
	return "guest_user"
}

// GET /api/wishlist
func (wc *WishlistController) GetWishlist(c *gin.Context) {
	userId := getUserIdFromContext(c)

	if config.DB != nil {
		var items []models.WishlistItem
		if err := config.DB.Preload("Product").Where("user_id = ?", userId).Find(&items).Error; err == nil {
			products := make([]models.Product, 0)
			for _, item := range items {
				if item.Product.ID != "" {
					products = append(products, item.Product)
				}
			}
			c.JSON(http.StatusOK, products)
			return
		}
	}

	// Memory Fallback
	mockWishlistsMutex.RLock()
	productIDs := mockWishlists[userId]
	mockWishlistsMutex.RUnlock()

	var result []models.Product
	for _, pid := range productIDs {
		for _, p := range mockProducts {
			if p.ID == pid {
				result = append(result, p)
				break
			}
		}
	}
	if result == nil {
		result = []models.Product{}
	}

	c.JSON(http.StatusOK, result)
}

// POST /api/wishlist/toggle
func (wc *WishlistController) ToggleWishlist(c *gin.Context) {
	var req models.ToggleWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := getUserIdFromContext(c)

	if config.DB != nil {
		var existing models.WishlistItem
		err := config.DB.Where("user_id = ? AND product_id = ?", userId, req.ProductID).First(&existing).Error
		if err == nil {
			// Item exists, remove it
			config.DB.Delete(&existing)
			c.JSON(http.StatusOK, gin.H{
				"status":     "removed",
				"isFavorite": false,
				"productId":  req.ProductID,
			})
			return
		}

		// Item does not exist, add it
		newItem := models.WishlistItem{
			ID:        fmt.Sprintf("wsh_%d", time.Now().UnixNano()),
			UserID:    userId,
			ProductID: req.ProductID,
			CreatedAt: time.Now(),
		}
		if err := config.DB.Create(&newItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wishlist"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     "added",
			"isFavorite": true,
			"productId":  req.ProductID,
		})
		return
	}

	// Memory Fallback
	mockWishlistsMutex.Lock()
	pIDs := mockWishlists[userId]
	found := false
	var updated []string
	for _, id := range pIDs {
		if id == req.ProductID {
			found = true
		} else {
			updated = append(updated, id)
		}
	}

	isFav := false
	statusStr := "removed"
	if !found {
		updated = append(updated, req.ProductID)
		isFav = true
		statusStr = "added"
	}
	mockWishlists[userId] = updated
	mockWishlistsMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"status":     statusStr,
		"isFavorite": isFav,
		"productId":  req.ProductID,
	})
}

// DELETE /api/wishlist/:productId
func (wc *WishlistController) RemoveWishlistItem(c *gin.Context) {
	productID := c.Param("productId")
	userId := getUserIdFromContext(c)

	if config.DB != nil {
		config.DB.Where("user_id = ? AND product_id = ?", userId, productID).Delete(&models.WishlistItem{})
		c.JSON(http.StatusOK, gin.H{
			"status":    "removed",
			"productId": productID,
		})
		return
	}

	// Memory Fallback
	mockWishlistsMutex.Lock()
	pIDs := mockWishlists[userId]
	var updated []string
	for _, id := range pIDs {
		if id != productID {
			updated = append(updated, id)
		}
	}
	mockWishlists[userId] = updated
	mockWishlistsMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"status":    "removed",
		"productId": productID,
	})
}
