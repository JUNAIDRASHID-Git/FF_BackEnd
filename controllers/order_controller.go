package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"funfillers/backend/config"
	"funfillers/backend/models"
	"github.com/gin-gonic/gin"
)

type OrderController struct{}

func NewOrderController() *OrderController {
	return &OrderController{}
}

var mockOrders = []models.Order{
	{
		ID:              "ord_1001",
		UserID:          "usr_demo",
		UserName:        "Demo Customer",
		UserEmail:       "user@funfillers.com",
		Items:           []models.OrderItem{{ProductID: "prod_1", ProductName: "FunFillers Premium Wireless Gaming Headset", Price: 129.99, Quantity: 1, ImageURL: "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=600"}},
		TotalAmount:     129.99,
		ShippingAddress: "123 Commerce St, Tech City, CA 94016",
		Status:          "Processing",
		PaymentStatus:   "Paid",
		CreatedAt:       time.Now().Add(-24 * time.Hour),
	},
	{
		ID:              "ord_1002",
		UserID:          "usr_demo",
		UserName:        "Demo Customer",
		UserEmail:       "user@funfillers.com",
		Items:           []models.OrderItem{{ProductID: "prod_2", ProductName: "Ergonomic RGB Mechanical Keyboard", Price: 89.95, Quantity: 1, ImageURL: "https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=600"}},
		TotalAmount:     89.95,
		ShippingAddress: "123 Commerce St, Tech City, CA 94016",
		Status:          "Delivered",
		PaymentStatus:   "Paid",
		CreatedAt:       time.Now().Add(-72 * time.Hour),
	},
}

func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, _ := c.Get("userId")
	userEmail, _ := c.Get("userEmail")

	userIdStr := fmt.Sprintf("%v", userId)
	userEmailStr := fmt.Sprintf("%v", userEmail)
	orderId := fmt.Sprintf("ord_%d", time.Now().UnixNano()/1e6)

	// Fetch customer real name from database if available
	userNameStr := "Valued Customer"
	if config.DB != nil && (userIdStr != "" || userEmailStr != "") {
		var u models.User
		if err := config.DB.Where("id = ? OR LOWER(email) = ?", userIdStr, strings.ToLower(userEmailStr)).First(&u).Error; err == nil && u.Name != "" {
			userNameStr = u.Name
		}
	}

	var total float64
	items := make([]models.OrderItem, len(req.Items))
	for i := range req.Items {
		total += req.Items[i].Price * float64(req.Items[i].Quantity)

		prodName := req.Items[i].ProductName
		imgUrl := req.Items[i].ImageURL

		// Resolve missing product name / image from database if not passed
		if (prodName == "" || imgUrl == "") && config.DB != nil && req.Items[i].ProductID != "" {
			var p models.Product
			if err := config.DB.Where("id = ?", req.Items[i].ProductID).First(&p).Error; err == nil {
				if prodName == "" {
					prodName = p.Name
				}
				if imgUrl == "" {
					imgUrl = p.ImageURL
				}
			}
		}

		items[i] = models.OrderItem{
			ID:          fmt.Sprintf("item_%d_%d", time.Now().UnixNano()/1e6, i),
			OrderID:     orderId,
			ProductID:   req.Items[i].ProductID,
			ProductName: prodName,
			Price:       req.Items[i].Price,
			Quantity:    req.Items[i].Quantity,
			ImageURL:    imgUrl,
		}
	}

	order := models.Order{
		ID:              orderId,
		UserID:          userIdStr,
		UserName:        userNameStr,
		UserEmail:       userEmailStr,
		Items:           items,
		TotalAmount:     total,
		ShippingAddress: req.ShippingAddress,
		Status:          "Pending",
		PaymentStatus:   "Paid",
		CreatedAt:       time.Now(),
	}

	if config.DB != nil {
		if err := config.DB.Create(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to persist order: " + err.Error()})
			return
		}
	}
	mockOrders = append([]models.Order{order}, mockOrders...)
	saveOrdersToFile()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order placed successfully",
		"order":   order,
	})
}

func (oc *OrderController) GetUserOrders(c *gin.Context) {
	userId, _ := c.Get("userId")
	userEmail, _ := c.Get("userEmail")

	userIdStr := fmt.Sprintf("%v", userId)
	userEmailStr := fmt.Sprintf("%v", userEmail)

	if config.DB != nil {
		var userOrders []models.Order
		query := config.DB.Preload("Items").Order("created_at desc")
		if userEmailStr != "" && userEmailStr != "<nil>" && userIdStr != "" && userIdStr != "<nil>" {
			query = query.Where("user_id = ? OR LOWER(user_email) = ?", userIdStr, strings.ToLower(userEmailStr))
		} else if userIdStr != "" && userIdStr != "<nil>" {
			query = query.Where("user_id = ?", userIdStr)
		} else if userEmailStr != "" && userEmailStr != "<nil>" {
			query = query.Where("LOWER(user_email) = ?", strings.ToLower(userEmailStr))
		}

		if err := query.Find(&userOrders).Error; err == nil {
			c.JSON(http.StatusOK, userOrders)
			return
		}
	}

	var userOrders []models.Order
	for _, o := range mockOrders {
		if o.UserID == userIdStr || (userEmailStr != "" && strings.EqualFold(o.UserEmail, userEmailStr)) {
			userOrders = append(userOrders, o)
		}
	}

	if userOrders == nil {
		userOrders = []models.Order{}
	}

	c.JSON(http.StatusOK, userOrders)
}
