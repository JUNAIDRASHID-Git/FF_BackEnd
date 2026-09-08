package models

import "time"

type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role      string    `json:"role" gorm:"default:'customer'"` // 'admin' or 'customer'
	Phone     string    `json:"phone,omitempty"`
	AvatarURL string    `json:"avatarUrl,omitempty"`
	Status    string    `json:"status" gorm:"default:'Active'"` // 'Active', 'Suspended'
	CreatedAt time.Time `json:"createdAt"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

type GoogleAuthRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl"`
	GoogleID  string `json:"googleId"`
}

type SendOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

