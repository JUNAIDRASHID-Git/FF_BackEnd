package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	JWTSecret         string
	Env               string
	RazorpayKeyID     string
	RazorpayKeySecret string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ No .env file found, using system environment variables or defaults.")
	} else {
		log.Println("🔒 Loaded environment configuration from .env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5050"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "funfillers_super_secret_jwt_key_2026"
	}

	env := os.Getenv("NODE_ENV")
	if env == "" {
		env = "development"
	}

	razorpayKeyID := os.Getenv("RAZORPAY_KEY_ID")
	if razorpayKeyID == "" {
		razorpayKeyID = "rzp_test_TaFCXrgrWjxR41"
	}

	razorpayKeySecret := os.Getenv("RAZORPAY_KEY_SECRET")
	if razorpayKeySecret == "" {
		razorpayKeySecret = "gqjcm0rxZ3AbrbQwarH2wZql"
	}

	return Config{
		Port:              port,
		JWTSecret:         jwtSecret,
		Env:               env,
		RazorpayKeyID:     razorpayKeyID,
		RazorpayKeySecret: razorpayKeySecret,
	}
}
