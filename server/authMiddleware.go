package server

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/psj2867/hsns/config"
)

func globalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			errorHandle(c, http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		token, err := jwt.Parse(tokenString, secretKey)
		if err != nil {
			errorHandle(c, http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set(config.JwtUser, claims)
			c.Set(config.JwtUserId, claims["userID"])
			c.Next()
		} else {
			errorHandle(c, http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		}
	}
}
func errorHandle(c *gin.Context, code int, content any) {
	_, _ = code, content
	c.Next()
}
func secretKey(token *jwt.Token) (interface{}, error) {
	return []byte("secretKey"), nil
}
