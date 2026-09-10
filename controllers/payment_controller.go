package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"time"

	"funfillers/backend/config"

	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

type PaymentController struct {
	cfg config.Config
}

func NewPaymentController(cfg config.Config) *PaymentController {
	return &PaymentController{cfg: cfg}
}

type CreateRazorpayOrderRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Receipt  string  `json:"receipt"`
}

type VerifyPaymentRequest struct {
	RazorpayOrderID   string `json:"razorpay_order_id"`
	RazorpayPaymentID string `json:"razorpay_payment_id"`
	RazorpaySignature string `json:"razorpay_signature"`
}

// POST /api/payment/create-order or /api/create-order
func (pc *PaymentController) CreateOrder(c *gin.Context) {
	var req CreateRazorpayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currency := req.Currency
	if currency == "" {
		currency = "INR"
	}

	// Razorpay expects amount in smallest currency unit (paise)
	amountInPaise := int(math.Round(req.Amount * 100))
	if amountInPaise < 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be at least ₹1 (100 paise)"})
		return
	}

	receipt := req.Receipt
	if receipt == "" {
		receipt = fmt.Sprintf("rcpt_%d", time.Now().UnixNano())
	}

	keyID := pc.cfg.RazorpayKeyID
	keySecret := pc.cfg.RazorpayKeySecret

	client := razorpay.NewClient(keyID, keySecret)

	data := map[string]interface{}{
		"amount":   amountInPaise,
		"currency": currency,
		"receipt":  receipt,
		"payment": map[string]interface{}{
			"capture": "automatic",
		},
	}

	orderData, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create Razorpay order",
			"details": err.Error(),
		})
		return
	}

	orderID, _ := orderData["id"].(string)

	c.JSON(http.StatusOK, gin.H{
		"order_id": orderID,
		"amount":   amountInPaise,
		"currency": currency,
		"key_id":   keyID,
	})
}

// POST /api/payment/verify or /api/verify-payment
func (pc *PaymentController) VerifyPayment(c *gin.Context) {
	var req VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RazorpayOrderID == "" || req.RazorpayPaymentID == "" || req.RazorpaySignature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required signature parameters"})
		return
	}

	keySecret := pc.cfg.RazorpayKeySecret
	data := req.RazorpayOrderID + "|" + req.RazorpayPaymentID

	h := hmac.New(sha256.New, []byte(keySecret))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(expectedSignature), []byte(req.RazorpaySignature)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failure",
			"error":  "Razorpay payment signature mismatch",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "Payment verified successfully",
		"order_id":   req.RazorpayOrderID,
		"payment_id": req.RazorpayPaymentID,
	})
}
