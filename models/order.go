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
	ID                 string      `json:"id" gorm:"primaryKey"`
	UserID             string      `json:"userId" gorm:"index"`
	UserName           string      `json:"userName"`
	UserEmail          string      `json:"userEmail"`
	TotalAmount        float64     `json:"totalAmount"`
	Status             string      `json:"status"` // 'Pending', 'Processing', 'Shipped', 'Delivered', 'Cancelled'
	PaymentMethod      string      `json:"paymentMethod"`
	PaymentID          string      `json:"paymentId"`
	RazorpayOrderID    string      `json:"razorpayOrderId"`
	PaymentStatus      string      `json:"paymentStatus"`
	ShippingAddress    string      `json:"shippingAddress"`
	RefundID           string      `json:"refundId"`
	RefundAmount       float64     `json:"refundAmount"`
	RefundStatus       string      `json:"refundStatus"`
	CancellationReason string      `json:"cancellationReason"`
	CancelledAt        *time.Time  `json:"cancelledAt"`
	Items              []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt          time.Time   `json:"createdAt"`
}

type CreateOrderRequest struct {
	Items           []OrderItem `json:"items" binding:"required"`
	ShippingAddress string      `json:"shippingAddress" binding:"required"`
	UserID          string      `json:"userId"`
	UserEmail       string      `json:"userEmail"`
	UserName        string      `json:"userName"`
	PaymentMethod   string      `json:"paymentMethod"`
	PaymentID       string      `json:"paymentId"`
	RazorpayOrderID string      `json:"razorpayOrderId"`
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
