package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/psj2867/hsns/config"
	"github.com/psj2867/hsns/models"
)

const Auth_Header = "Authorization"

type AuthUserType = map[string]any

func GlobalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(Auth_Header)
		if tokenString == "" {
			errorHandle(c, http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		token, err := ParseJwt(tokenString)
		if err != nil {
			errorHandle(c, http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set(config.JwtUser, AuthUserType(claims))
			c.Next()
		} else {
			c.Set(config.JwtUser, AuthUserType{})
			errorHandle(c, http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		}
	}
}
func ParseJwt(tokenString string, options ...jwt.ParserOption) (*jwt.Token, error) {
	return jwt.Parse(tokenString, secretKey, options...)
}
func errorHandle(c *gin.Context, code int, content any) {
	_, _ = code, content
	c.Next()
}
func getSecretKey() []byte {
	return []byte("secretkey")
}
func secretKey(token *jwt.Token) (interface{}, error) {
	return getSecretKey(), nil
}
func GetAuthInfo(c *gin.Context) (AuthUserType, bool) {
	a, ok := c.Get(config.JwtUser)
	if !ok {
		return nil, false
	}
	return a.(AuthUserType), true
}
func GetAuthInfoByKey(c *gin.Context, key string) (any, bool) {
	claims, ok := c.Get(config.JwtUser)
	if !ok {
		return nil, false
	}
	val := claims.(AuthUserType)[key]
	if val == nil {
		return nil, false
	}
	return val, true
}
func MustGetAuthInfoByKey(c *gin.Context, key string) any {
	val, ok := GetAuthInfoByKey(c, key)
	if !ok {
		panic("should use shouldAuth middleware")
	}
	return val
}
func HasAuth(c *gin.Context) bool {
	_, ok := GetAuthInfo(c)
	return ok
}

func generateToken(values map[string]any) (string, error) {
	accessTokenClaims := jwt.MapClaims(values)
	accessTokenClaims["iat"] = time.Now().Unix()
	accessTokenClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken, err := token.SignedString(getSecretKey())
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func GenerateToken(user *models.User) (string, error) {
	return generateToken(map[string]any{
		"userId": user.UserId,
	})
}
