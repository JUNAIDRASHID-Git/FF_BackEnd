package models

import "time"

type OrderItem struct {
	ID          string  `json:"id" gorm:"primaryKey"`
	OrderID     string  `json:"orderId" gorm:"index"`
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	ImageURL    string  `json:"imageUrl"`
}

type Order struct {
	ID              string      `json:"id" gorm:"primaryKey"`
	UserID          string      `json:"userId" gorm:"index"`
	UserName        string      `json:"userName"`
	UserEmail       string      `json:"userEmail"`
	TotalAmount     float64     `json:"totalAmount"`
	Status          string      `json:"status"` // 'Pending', 'Processing', 'Shipped', 'Delivered', 'Cancelled'
	PaymentStatus   string      `json:"paymentStatus"`
	ShippingAddress string      `json:"shippingAddress"`
	Items           []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt       time.Time   `json:"createdAt"`
}

type CreateOrderRequest struct {
	Items           []OrderItem `json:"items" binding:"required"`
	ShippingAddress string      `json:"shippingAddress" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AdminStats struct {
	TotalRevenue   float64 `json:"totalRevenue"`
	TotalOrders    int     `json:"totalOrders"`
	TotalProducts  int     `json:"totalProducts"`
	TotalCustomers int     `json:"totalCustomers"`
}
