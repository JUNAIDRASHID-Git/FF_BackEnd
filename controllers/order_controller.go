package controllers

import (
	"fmt"
	"net/http"
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

	var total float64
	for i := range req.Items {
		total += req.Items[i].Price * float64(req.Items[i].Quantity)
		if req.Items[i].ProductName == "" || req.Items[i].ImageURL == "" {
			for _, p := range mockProducts {
				if p.ID == req.Items[i].ProductID {
					if req.Items[i].ProductName == "" {
						req.Items[i].ProductName = p.Name
					}
					if req.Items[i].ImageURL == "" {
						req.Items[i].ImageURL = p.ImageURL
					}
					break
				}
			}
		}
	}

	order := models.Order{
		ID:              fmt.Sprintf("ord_%d", time.Now().UnixNano()/1e6),
		UserID:          fmt.Sprintf("%v", userId),
		UserName:        "Valued Customer",
		UserEmail:       fmt.Sprintf("%v", userEmail),
		Items:           req.Items,
		TotalAmount:     total,
		ShippingAddress: req.ShippingAddress,
		Status:          "Pending",
		PaymentStatus:   "Paid",
		CreatedAt:       time.Now(),
	}

	if config.DB != nil {
		config.DB.Create(&order)
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

	if config.DB != nil {
		var userOrders []models.Order
		if err := config.DB.Preload("Items").Where("user_id = ?", fmt.Sprintf("%v", userId)).Find(&userOrders).Error; err == nil {
			c.JSON(http.StatusOK, userOrders)
			return
		}
	}

	var userOrders []models.Order
	for _, o := range mockOrders {
		if o.UserID == fmt.Sprintf("%v", userId) {
			userOrders = append(userOrders, o)
		}
	}

	if userOrders == nil {
		userOrders = []models.Order{}
	}

	c.JSON(http.StatusOK, userOrders)
}
