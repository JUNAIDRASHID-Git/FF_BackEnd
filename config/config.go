package config

import "os"

type Config struct {
	Port      string
	JWTSecret string
	Env       string
}

func LoadConfig() Config {
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

	return Config{
		Port:      port,
		JWTSecret: jwtSecret,
		Env:       env,
	}
}
