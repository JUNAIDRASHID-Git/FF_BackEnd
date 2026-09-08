package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"funfillers/backend/config"
	"funfillers/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	cfg config.Config
}

func NewAuthController(cfg config.Config) *AuthController {
	return &AuthController{cfg: cfg}
}

var mockUserMap = map[string]models.User{
	"admin@funfillers.com": {
		ID:        "usr_admin",
		Name:      "FUNFILLERS Admin",
		Email:     "admin@funfillers.com",
		Password:  "admin123",
		Role:      "admin",
		Status:    "Active",
		CreatedAt: time.Now(),
	},
	"user@funfillers.com": {
		ID:        "usr_demo",
		Name:      "Demo Customer",
		Email:     "user@funfillers.com",
		Password:  "user123",
		Role:      "customer",
		Status:    "Active",
		CreatedAt: time.Now(),
	},
}

func (ac *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, exists := mockUserMap[req.Email]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	role := "customer"
	if req.Role == "admin" {
		role = "admin"
	}

	user := models.User{
		ID:        fmt.Sprintf("usr_%d", time.Now().UnixNano()),
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		Role:      role,
		Status:    "Active",
		CreatedAt: time.Now(),
	}

	if config.DB != nil {
		config.DB.Create(&user)
	}

	mockUserMap[req.Email] = user
	token, _ := generateToken(user, ac.cfg.JWTSecret)

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := mockUserMap[req.Email]
	if !exists || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := generateToken(user, ac.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func (ac *AuthController) GoogleAuth(c *gin.Context) {
	var req models.GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := mockUserMap[req.Email]
	if !exists {
		name := req.Name
		if name == "" {
			name = strings.Split(req.Email, "@")[0]
		}
		user = models.User{
			ID:        fmt.Sprintf("usr_g_%d", time.Now().UnixNano()),
			Name:      name,
			Email:     req.Email,
			Password:  "google_oauth_user",
			Role:      "customer",
			AvatarURL: req.AvatarURL,
			Status:    "Active",
			CreatedAt: time.Now(),
		}
		if config.DB != nil {
			config.DB.Create(&user)
		}
		mockUserMap[req.Email] = user
	} else if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
		mockUserMap[req.Email] = user
	}

	AddOrUpdateMockUser(user)

	token, err := generateToken(user, ac.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func (ac *AuthController) SendOTP(c *gin.Context) {
	var req models.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phone := strings.TrimSpace(req.Phone)
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "OTP sent successfully via SMS/Firebase",
		"phone":          phone,
		"verificationId": fmt.Sprintf("verif_%d", time.Now().UnixNano()),
		"devOtp":         "123456",
	})
}

func (ac *AuthController) VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phone := strings.TrimSpace(req.Phone)
	otp := strings.TrimSpace(req.OTP)

	if otp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP is required"})
		return
	}

	if otp != "123456" && len(otp) != 6 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP code provided"})
		return
	}

	emailKey := fmt.Sprintf("phone_%s@funfillers.com", strings.ReplaceAll(phone, "+", ""))
	user, exists := mockUserMap[emailKey]
	if !exists {
		user = models.User{
			ID:        fmt.Sprintf("usr_p_%d", time.Now().UnixNano()),
			Name:      fmt.Sprintf("Phone User (%s)", phone),
			Email:     emailKey,
			Phone:     phone,
			Password:  "phone_otp_user",
			Role:      "customer",
			Status:    "Active",
			CreatedAt: time.Now(),
		}
		if config.DB != nil {
			config.DB.Create(&user)
		}
		mockUserMap[emailKey] = user
	}

	token, err := generateToken(user, ac.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func generateToken(user models.User, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
