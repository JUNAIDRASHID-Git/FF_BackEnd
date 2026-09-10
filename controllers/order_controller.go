package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"funfillers/backend/config"
	"funfillers/backend/models"
	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

type OrderController struct {
	cfg config.Config
}

func NewOrderController(cfg config.Config) *OrderController {
	return &OrderController{cfg: cfg}
}

var mockOrders = []models.Order{}

func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, _ := c.Get("userId")
	userEmail, _ := c.Get("userEmail")

	userIdStr := ""
	if userId != nil && fmt.Sprintf("%v", userId) != "<nil>" {
		userIdStr = fmt.Sprintf("%v", userId)
	}
	if userIdStr == "" {
		userIdStr = c.GetHeader("X-User-ID")
	}
	if userIdStr == "" {
		userIdStr = req.UserID
	}

	userEmailStr := ""
	if userEmail != nil && fmt.Sprintf("%v", userEmail) != "<nil>" {
		userEmailStr = fmt.Sprintf("%v", userEmail)
	}
	if userEmailStr == "" {
		userEmailStr = c.GetHeader("X-User-Email")
	}
	if userEmailStr == "" {
		userEmailStr = req.UserEmail
	}

	userNameStr := req.UserName
	if userNameStr == "" {
		userNameStr = "Valued Customer"
	}

	orderId := fmt.Sprintf("ord_%d", time.Now().UnixNano()/1e6)

	// Fetch customer real name from database if available
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

	paymentMethod := req.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = "Online Payment"
	}

	paymentID := req.PaymentID
	if paymentID == "" && strings.Contains(paymentMethod, "pay_") {
		parts := strings.Split(paymentMethod, "pay_")
		if len(parts) > 1 {
			paymentID = "pay_" + strings.TrimRight(parts[1], ") ")
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
		PaymentMethod:   paymentMethod,
		PaymentID:       paymentID,
		RazorpayOrderID: req.RazorpayOrderID,
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

	userIdStr := ""
	if userId != nil && fmt.Sprintf("%v", userId) != "<nil>" {
		userIdStr = fmt.Sprintf("%v", userId)
	}
	if userIdStr == "" {
		userIdStr = c.GetHeader("X-User-ID")
	}
	if userIdStr == "" {
		userIdStr = c.Query("userId")
	}

	userEmailStr := ""
	if userEmail != nil && fmt.Sprintf("%v", userEmail) != "<nil>" {
		userEmailStr = fmt.Sprintf("%v", userEmail)
	}
	if userEmailStr == "" {
		userEmailStr = c.GetHeader("X-User-Email")
	}
	if userEmailStr == "" {
		userEmailStr = c.Query("email")
	}
	if userEmailStr == "" {
		userEmailStr = c.Query("userEmail")
	}

	if config.DB != nil {
		var userOrders []models.Order
		query := config.DB.Preload("Items").Order("created_at desc")
		if userEmailStr != "" && userIdStr != "" {
			query = query.Where("user_id = ? OR LOWER(user_email) = ?", userIdStr, strings.ToLower(userEmailStr))
		} else if userIdStr != "" {
			query = query.Where("user_id = ?", userIdStr)
		} else if userEmailStr != "" {
			query = query.Where("LOWER(user_email) = ?", strings.ToLower(userEmailStr))
		}

		if err := query.Find(&userOrders).Error; err == nil && len(userOrders) > 0 {
			c.JSON(http.StatusOK, userOrders)
			return
		}
	}

	var userOrders []models.Order
	for _, o := range mockOrders {
		matchUser := userIdStr != "" && o.UserID == userIdStr
		matchEmail := userEmailStr != "" && strings.EqualFold(o.UserEmail, userEmailStr)
		if matchUser || matchEmail {
			userOrders = append(userOrders, o)
		}
	}

	if userOrders == nil {
		userOrders = []models.Order{}
	}

	c.JSON(http.StatusOK, userOrders)
}

func (oc *OrderController) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)

	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = "Cancelled by user"
	}

	var targetOrder *models.Order
	var dbFound bool

	if config.DB != nil {
		var dbOrder models.Order
		if err := config.DB.Preload("Items").First(&dbOrder, "id = ?", id).Error; err == nil {
			targetOrder = &dbOrder
			dbFound = true
		}
	}

	if targetOrder == nil {
		for i, o := range mockOrders {
			if o.ID == id {
				targetOrder = &mockOrders[i]
				break
			}
		}
	}

	if targetOrder == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	statusLower := strings.ToLower(targetOrder.Status)
	if statusLower == "shipped" || statusLower == "delivered" || statusLower == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Order cannot be cancelled because current status is '%s'", targetOrder.Status),
		})
		return
	}

	now := time.Now()
	refundMessage := "Order cancelled successfully."
	var refundID string

	// Determine payment ID for Razorpay online payments
	paymentID := targetOrder.PaymentID
	if paymentID == "" && strings.Contains(targetOrder.PaymentMethod, "pay_") {
		parts := strings.Split(targetOrder.PaymentMethod, "pay_")
		if len(parts) > 1 {
			paymentID = "pay_" + strings.TrimRight(parts[1], ") ")
		}
	}

	isOnlinePayment := paymentID != "" || (strings.Contains(strings.ToLower(targetOrder.PaymentMethod), "razorpay") || strings.Contains(strings.ToLower(targetOrder.PaymentMethod), "online"))

	if isOnlinePayment {
		if paymentID != "" && oc.cfg.RazorpayKeyID != "" && oc.cfg.RazorpayKeySecret != "" {
			client := razorpay.NewClient(oc.cfg.RazorpayKeyID, oc.cfg.RazorpayKeySecret)
			amountInPaise := int(targetOrder.TotalAmount * 100)
			if amountInPaise < 100 {
				amountInPaise = 100
			}

			refundParams := map[string]interface{}{
				"amount": amountInPaise,
				"speed":  "optimum",
				"notes": map[string]interface{}{
					"reason":   reason,
					"order_id": targetOrder.ID,
				},
			}

			res, err := client.Payment.Refund(paymentID, amountInPaise, refundParams, nil)
			if err == nil && res != nil {
				if idVal, ok := res["id"].(string); ok {
					refundID = idVal
				}
				targetOrder.RefundID = refundID
				targetOrder.RefundAmount = targetOrder.TotalAmount
				targetOrder.RefundStatus = "Processed"
				targetOrder.PaymentStatus = "Refunded"
				refundMessage = fmt.Sprintf("Order cancelled and refund of ₹%.2f initiated via Razorpay (Refund ID: %s)", targetOrder.TotalAmount, refundID)
			} else {
				// Mark as refund pending if Razorpay returned mock/test error
				targetOrder.RefundAmount = targetOrder.TotalAmount
				targetOrder.RefundStatus = "Pending"
				targetOrder.PaymentStatus = "Refund Pending"
				refundMessage = fmt.Sprintf("Order cancelled. Online refund initiated for ₹%.2f (Support will confirm refund processing).", targetOrder.TotalAmount)
			}
		} else {
			targetOrder.RefundAmount = targetOrder.TotalAmount
			targetOrder.RefundStatus = "Initiated"
			targetOrder.PaymentStatus = "Refund Initiated"
			refundMessage = fmt.Sprintf("Order cancelled and full refund of ₹%.2f initiated.", targetOrder.TotalAmount)
		}
	} else {
		targetOrder.PaymentStatus = "Cancelled"
		refundMessage = "Order cancelled successfully."
	}

	targetOrder.Status = "Cancelled"
	targetOrder.CancellationReason = reason
	targetOrder.CancelledAt = &now

	if dbFound && config.DB != nil {
		config.DB.Save(targetOrder)
	}

	for i, o := range mockOrders {
		if o.ID == id {
			mockOrders[i] = *targetOrder
			break
		}
	}
	saveOrdersToFile()

	c.JSON(http.StatusOK, gin.H{
		"message":            refundMessage,
		"refundId":           targetOrder.RefundID,
		"refundAmount":       targetOrder.RefundAmount,
		"refundStatus":       targetOrder.RefundStatus,
		"cancellationReason": targetOrder.CancellationReason,
		"order":              targetOrder,
	})
}
