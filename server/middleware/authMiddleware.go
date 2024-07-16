package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const Auth_Header = "Authorization"
const TokenKey_UserId = "userid"

const (
	MiddlewareJwtUser   = "server/middleware/user"
	MiddlewareJwtUserId = "server/middleware/userId"
)

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
			c.Set(MiddlewareJwtUser, AuthUserType(claims))
			c.Next()
		} else {
			c.Set(MiddlewareJwtUser, AuthUserType{})
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
func GetAuthUserId(c *gin.Context) (string, bool) {
	a, ok := GetAuthInfoByKey(c, TokenKey_UserId)
	if !ok {
		return "", false
	}
	b, ok := a.(string)
	if !ok {
		return "", false
	}
	return b, true

}
func GetAuthInfo(c *gin.Context) (AuthUserType, bool) {
	a, ok := c.Get(MiddlewareJwtUser)
	if !ok {
		return nil, false
	}
	return a.(AuthUserType), true
}
func GetAuthInfoByKey(c *gin.Context, key string) (any, bool) {
	claims, ok := c.Get(MiddlewareJwtUser)
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

type UserInfoForToken struct {
	UserId string
	values map[string]any
}

func GenerateToken(userInfo UserInfoForToken, options ...func(UserInfoForToken)) (string, error) {
	userInfo.values = map[string]any{}
	for _, v := range options {
		v(userInfo)
	}
	accessTokenClaims := jwt.MapClaims(userInfo.values)
	accessTokenClaims[TokenKey_UserId] = userInfo.UserId
	accessTokenClaims["iat"] = time.Now().Unix()
	accessTokenClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken, err := token.SignedString(getSecretKey())
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
