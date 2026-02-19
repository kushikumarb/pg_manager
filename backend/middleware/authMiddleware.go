package middleware

import (
	"fmt"
	"net/http"
	"os" // Added for os.Getenv
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// 2. Extract the token (Format: "Bearer <token>")
		// Using strings.Fields is safer as it handles extra spaces automatically
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 3. Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is what we expect (HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// MATCH YOUR SERVICE LOGIC: Get secret key from .env or fallback
			secretKey := os.Getenv("JWT_SECRET")
			if secretKey == "" {
				secretKey = "my_secret_key"
			}
			return []byte(secretKey), nil
		})

		// 4. Check if token is valid
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check Expiration
			if exp, ok := claims["exp"].(float64); ok {
				if float64(time.Now().Unix()) > exp {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
					c.Abort()
					return
				}
			}

			// 5. Attach User info to the context
			// In JWT, numbers are parsed as float64.
			// Handlers can convert this back to uint using: uint(val.(float64))
			c.Set("user_id", claims["user_id"])
			c.Set("role", claims["role"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			c.Abort()
			return
		}
	}
}
