package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"funfillers/backend/config"
	"funfillers/backend/models"
	"github.com/gin-gonic/gin"
)

type AdminController struct{}

const ordersFilePath = "./data/orders.json"
const usersFilePath = "./data/users.json"

func loadOrdersFromFile() {
	if _, err := os.Stat(ordersFilePath); os.IsNotExist(err) {
		saveOrdersToFile()
		return
	}
	data, err := os.ReadFile(ordersFilePath)
	if err != nil {
		return
	}
	var loaded []models.Order
	if err := json.Unmarshal(data, &loaded); err == nil && len(loaded) > 0 {
		for i := range loaded {
			for j := range loaded[i].Items {
				if loaded[i].Items[j].ProductName == "" || loaded[i].Items[j].ImageURL == "" {
					for _, p := range mockProducts {
						if p.ID == loaded[i].Items[j].ProductID {
							if loaded[i].Items[j].ProductName == "" {
								loaded[i].Items[j].ProductName = p.Name
							}
							if loaded[i].Items[j].ImageURL == "" {
								loaded[i].Items[j].ImageURL = p.ImageURL
							}
							break
						}
					}
				}
			}
		}
		mockOrders = loaded
	}
}

func saveOrdersToFile() {
	os.MkdirAll("./data", 0755)
	data, err := json.MarshalIndent(mockOrders, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(ordersFilePath, data, 0644)
}

func loadUsersFromFile() {
	if _, err := os.Stat(usersFilePath); os.IsNotExist(err) {
		saveUsersToFile()
		return
	}
	data, err := os.ReadFile(usersFilePath)
	if err != nil {
		return
	}
	var loaded []models.User
	if err := json.Unmarshal(data, &loaded); err == nil && len(loaded) > 0 {
		mockUsers = loaded
	}
}

func saveUsersToFile() {
	os.MkdirAll("./data", 0755)
	data, err := json.MarshalIndent(mockUsers, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(usersFilePath, data, 0644)
}

func NewAdminController() *AdminController {
	loadOrdersFromFile()
	loadUsersFromFile()
	return &AdminController{}
}

func AddOrUpdateMockUser(u models.User) {
	found := false
	for i, existing := range mockUsers {
		if strings.EqualFold(existing.Email, u.Email) || (u.ID != "" && existing.ID == u.ID) {
			if u.Name != "" {
				mockUsers[i].Name = u.Name
			}
			if u.AvatarURL != "" {
				mockUsers[i].AvatarURL = u.AvatarURL
			}
			found = true
			break
		}
	}
	if !found {
		if u.Role == "" {
			u.Role = "customer"
		}
		if u.Status == "" {
			u.Status = "Active"
		}
		mockUsers = append(mockUsers, u)
	}
	saveUsersToFile()
}

var mockUsers = []models.User{
	{ID: "usr_super1", Name: "Junaid Rashid", Email: "junaidrashid678@gmail.com", Role: "super_admin", Status: "Active"},
}

func (ac *AdminController) GetStats(c *gin.Context) {
	var totalRev float64
	for _, o := range mockOrders {
		totalRev += o.TotalAmount
	}

	stats := models.AdminStats{
		TotalRevenue:   totalRev,
		TotalOrders:    len(mockOrders),
		TotalProducts:  len(mockProducts),
		TotalCustomers: len(mockUsers),
	}

	if config.DB != nil {
		var rev float64
		var ordersCount, prodCount, userCount int64

		config.DB.Model(&models.Order{}).Select("COALESCE(SUM(total_amount), 0)").Scan(&rev)
		config.DB.Model(&models.Order{}).Count(&ordersCount)
		config.DB.Model(&models.Product{}).Count(&prodCount)
		config.DB.Model(&models.User{}).Count(&userCount)

		if ordersCount > 0 || prodCount > 0 {
			stats.TotalRevenue = rev
			stats.TotalOrders = int(ordersCount)
			stats.TotalProducts = int(prodCount)
			stats.TotalCustomers = int(userCount)
		}
	}

	c.JSON(http.StatusOK, stats)
}

func (ac *AdminController) GetUsers(c *gin.Context) {
	if config.DB != nil {
		var users []models.User
		if err := config.DB.Find(&users).Error; err == nil && len(users) > 0 {
			c.JSON(http.StatusOK, users)
			return
		}
	}
	c.JSON(http.StatusOK, mockUsers)
}

func isSuperAdminEmail(email string) bool {
	e := strings.ToLower(strings.TrimSpace(email))
	return e == "junaidrashid678@gmail.com" || e == "juanidrashid678@gmail.com"
}

func (ac *AdminController) UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, u := range mockUsers {
		if u.ID == id {
			if isSuperAdminEmail(u.Email) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Super Admin status cannot be altered"})
				return
			}
			mockUsers[i].Status = req.Status
			saveUsersToFile()

			c.JSON(http.StatusOK, gin.H{
				"message": "User status updated",
				"user":   mockUsers[i],
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func (ac *AdminController) UpdateUserRole(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, u := range mockUsers {
		if u.ID == id {
			if isSuperAdminEmail(u.Email) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Super Admin role cannot be altered"})
				return
			}
			mockUsers[i].Role = req.Role
			saveUsersToFile()

			c.JSON(http.StatusOK, gin.H{
				"message": "User role updated",
				"user":   mockUsers[i],
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func (ac *AdminController) RequestAdminAccess(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Name  string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailLower := strings.ToLower(strings.TrimSpace(req.Email))
	isSuper := isSuperAdminEmail(emailLower)

	for _, u := range mockUsers {
		if strings.ToLower(u.Email) == emailLower {
			c.JSON(http.StatusOK, gin.H{
				"status":       u.Status,
				"role":         u.Role,
				"isSuperAdmin": isSuper || u.Role == "super_admin",
				"user":         u,
			})
			return
		}
	}

	newStatus := "Pending"
	role := "admin"
	if isSuper {
		newStatus = "Active"
		role = "super_admin"
	}

	name := req.Name
	if name == "" {
		name = strings.Split(req.Email, "@")[0]
	}

	newUser := models.User{
		ID:        fmt.Sprintf("usr_%d", time.Now().UnixNano()/1e6),
		Name:      name,
		Email:     req.Email,
		Role:      role,
		Status:    newStatus,
		CreatedAt: time.Now(),
	}

	mockUsers = append([]models.User{newUser}, mockUsers...)
	saveUsersToFile()

	c.JSON(http.StatusCreated, gin.H{
		"status":       newUser.Status,
		"role":         newUser.Role,
		"isSuperAdmin": isSuper,
		"user":         newUser,
	})
}

func (ac *AdminController) GetAllOrders(c *gin.Context) {
	if config.DB != nil {
		var orders []models.Order
		if err := config.DB.Preload("Items").Order("created_at desc").Find(&orders).Error; err == nil && len(orders) > 0 {
			c.JSON(http.StatusOK, orders)
			return
		}
	}
	for i := range mockOrders {
		for j := range mockOrders[i].Items {
			if mockOrders[i].Items[j].ProductName == "" || mockOrders[i].Items[j].ImageURL == "" {
				for _, p := range mockProducts {
					if p.ID == mockOrders[i].Items[j].ProductID {
						if mockOrders[i].Items[j].ProductName == "" {
							mockOrders[i].Items[j].ProductName = p.Name
						}
						if mockOrders[i].Items[j].ImageURL == "" {
							mockOrders[i].Items[j].ImageURL = p.ImageURL
						}
						break
					}
				}
			}
		}
	}
	c.JSON(http.StatusOK, mockOrders)
}

func (ac *AdminController) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	normalizedStatus := req.Status
	switch strings.ToLower(req.Status) {
	case "pending":
		normalizedStatus = "Pending"
	case "processing":
		normalizedStatus = "Processing"
	case "shipped":
		normalizedStatus = "Shipped"
	case "delivered":
		normalizedStatus = "Delivered"
	case "cancelled":
		normalizedStatus = "Cancelled"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order status. Allowed: Pending, Processing, Shipped, Delivered, Cancelled"})
		return
	}

	if config.DB != nil {
		var dbOrder models.Order
		if err := config.DB.First(&dbOrder, "id = ?", id).Error; err == nil {
			dbOrder.Status = normalizedStatus
			config.DB.Save(&dbOrder)

			for i, o := range mockOrders {
				if o.ID == id {
					mockOrders[i].Status = normalizedStatus
					break
				}
			}
			saveOrdersToFile()

			c.JSON(http.StatusOK, gin.H{
				"message": "Order status updated",
				"order":   dbOrder,
			})
			return
		}
	}

	for i, o := range mockOrders {
		if o.ID == id {
			mockOrders[i].Status = normalizedStatus
			saveOrdersToFile()

			c.JSON(http.StatusOK, gin.H{
				"message": "Order status updated",
				"order":   mockOrders[i],
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
}

