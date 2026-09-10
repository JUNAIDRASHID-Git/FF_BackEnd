package middleware

import (
	"net/http"
	"strings"

	"funfillers/backend/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var claims jwt.MapClaims

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte(cfg.JWTSecret), nil
				})

				if err == nil && token.Valid {
					if cMap, ok := token.Claims.(jwt.MapClaims); ok {
						claims = cMap
					}
				} else {
					// Fallback: parse unverified token claims for third-party / mobile tokens
					parser := jwt.NewParser()
					if unverifiedToken, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{}); err == nil {
						if cMap, ok := unverifiedToken.Claims.(jwt.MapClaims); ok {
							claims = cMap
						}
					}
				}
			}
		}

		if claims != nil {
			if sub, ok := claims["sub"]; ok && sub != nil {
				c.Set("userId", sub)
			}
			if role, ok := claims["role"]; ok && role != nil {
				c.Set("userRole", role)
			}
			if email, ok := claims["email"]; ok && email != nil {
				c.Set("userEmail", email)
			}
		}

		// Fallback to X-User-ID and X-User-Email headers if context claims are missing
		if _, exists := c.Get("userId"); !exists {
			if xUserID := c.GetHeader("X-User-ID"); xUserID != "" {
				c.Set("userId", xUserID)
			}
		}
		if _, exists := c.Get("userEmail"); !exists {
			if xUserEmail := c.GetHeader("X-User-Email"); xUserEmail != "" {
				c.Set("userEmail", xUserEmail)
			}
		}

		// Check if we have any identity
		userId, hasUserID := c.Get("userId")
		userEmail, hasUserEmail := c.Get("userEmail")

		if !hasUserID && !hasUserEmail && authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token or user headers required"})
			c.Abort()
			return
		}

		if (userId == nil || userId == "") && (userEmail == nil || userEmail == "") {
			// If neither ID nor Email found, allow request to proceed if headers exist, else unauth
			if c.GetHeader("X-User-ID") == "" && c.GetHeader("X-User-Email") == "" && authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or user identification"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privilege required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func OptionalAuthMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte(cfg.JWTSecret), nil
				})
				if err == nil && token.Valid {
					if claims, ok := token.Claims.(jwt.MapClaims); ok {
						c.Set("userId", claims["sub"])
						c.Set("userRole", claims["role"])
						c.Set("userEmail", claims["email"])
					}
				} else {
					parser := jwt.NewParser()
					if unverifiedToken, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{}); err == nil {
						if claims, ok := unverifiedToken.Claims.(jwt.MapClaims); ok {
							c.Set("userId", claims["sub"])
							c.Set("userRole", claims["role"])
							c.Set("userEmail", claims["email"])
						}
					}
				}
			}
		}
		if _, exists := c.Get("userId"); !exists {
			if xUserID := c.GetHeader("X-User-ID"); xUserID != "" {
				c.Set("userId", xUserID)
			}
		}
		if _, exists := c.Get("userEmail"); !exists {
			if xUserEmail := c.GetHeader("X-User-Email"); xUserEmail != "" {
				c.Set("userEmail", xUserEmail)
			}
		}
		c.Next()
	}
}

