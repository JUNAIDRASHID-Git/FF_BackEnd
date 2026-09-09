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

type AddressController struct{}

func NewAddressController() *AddressController {
	return &AddressController{}
}

// In-memory fallback storage for addresses when DB is disabled
var (
	mockAddressesMutex sync.RWMutex
	mockAddresses      = make(map[string][]models.Address) // userId -> []Address
)

// Helper to extract userId from Auth token or headers or query
func getAddressUserId(c *gin.Context) string {
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

// GET /api/user/addresses
func (ac *AddressController) GetAddresses(c *gin.Context) {
	userId := getAddressUserId(c)

	if config.DB != nil {
		var addresses []models.Address
		if err := config.DB.Where("user_id = ?", userId).Order("is_default desc, created_at desc").Find(&addresses).Error; err == nil {
			if addresses == nil {
				addresses = []models.Address{}
			}
			c.JSON(http.StatusOK, addresses)
			return
		}
	}

	// In-memory fallback
	mockAddressesMutex.RLock()
	addresses := mockAddresses[userId]
	mockAddressesMutex.RUnlock()

	if addresses == nil {
		addresses = []models.Address{}
	}

	c.JSON(http.StatusOK, addresses)
}

// POST /api/user/addresses
func (ac *AddressController) SaveAddress(c *gin.Context) {
	var address models.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := getAddressUserId(c)
	address.UserID = userId

	if address.ID == "" {
		address.ID = fmt.Sprintf("addr_%d", time.Now().UnixNano())
	}
	if address.CreatedAt.IsZero() {
		address.CreatedAt = time.Now()
	}

	if config.DB != nil {
		// If new address is default, unset is_default for other addresses of this user
		if address.IsDefault {
			config.DB.Model(&models.Address{}).Where("user_id = ?", userId).Update("is_default", false)
		}

		// Save or update address
		if err := config.DB.Save(&address).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save address to database"})
			return
		}
		c.JSON(http.StatusOK, address)
		return
	}

	// In-memory fallback
	mockAddressesMutex.Lock()
	existingList := mockAddresses[userId]

	if address.IsDefault {
		for i := range existingList {
			existingList[i].IsDefault = false
		}
	} else if len(existingList) == 0 {
		address.IsDefault = true
	}

	// Check if updating existing
	updated := false
	for i, a := range existingList {
		if a.ID == address.ID {
			existingList[i] = address
			updated = true
			break
		}
	}
	if !updated {
		existingList = append([]models.Address{address}, existingList...)
	}

	mockAddresses[userId] = existingList
	mockAddressesMutex.Unlock()

	c.JSON(http.StatusOK, address)
}

// PUT /api/user/addresses/:id/default
func (ac *AddressController) SetDefaultAddress(c *gin.Context) {
	addressID := c.Param("id")
	userId := getAddressUserId(c)

	if config.DB != nil {
		config.DB.Model(&models.Address{}).Where("user_id = ?", userId).Update("is_default", false)
		config.DB.Model(&models.Address{}).Where("user_id = ? AND id = ?", userId, addressID).Update("is_default", true)

		var updatedAddress models.Address
		config.DB.Where("user_id = ? AND id = ?", userId, addressID).First(&updatedAddress)
		c.JSON(http.StatusOK, updatedAddress)
		return
	}

	// In-memory fallback
	mockAddressesMutex.Lock()
	existingList := mockAddresses[userId]
	var target models.Address
	for i := range existingList {
		if existingList[i].ID == addressID {
			existingList[i].IsDefault = true
			target = existingList[i]
		} else {
			existingList[i].IsDefault = false
		}
	}
	mockAddresses[userId] = existingList
	mockAddressesMutex.Unlock()

	c.JSON(http.StatusOK, target)
}

// DELETE /api/user/addresses/:id
func (ac *AddressController) DeleteAddress(c *gin.Context) {
	addressID := c.Param("id")
	userId := getAddressUserId(c)

	if config.DB != nil {
		config.DB.Where("user_id = ? AND id = ?", userId, addressID).Delete(&models.Address{})
		c.JSON(http.StatusOK, gin.H{"status": "deleted", "id": addressID})
		return
	}

	// In-memory fallback
	mockAddressesMutex.Lock()
	existingList := mockAddresses[userId]
	var updated []models.Address
	for _, a := range existingList {
		if a.ID != addressID {
			updated = append(updated, a)
		}
	}
	mockAddresses[userId] = updated
	mockAddressesMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "deleted", "id": addressID})
}
