package middlewares

import (
	"net/http"

	"github.com/Immortanbird/api_using_go/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	JWTSecretKey string
}

func (mw *Middleware) Authenticate(c *gin.Context) {
	access_token, err := c.Cookie("access_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	claims, err := utils.ParseToken(access_token, mw.JWTSecretKey)
	if err != nil {
		// Handle different types of errors
		if err == jwt.ErrTokenExpired {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		}
		return
	}

	c.Set("userID", claims.UserID)

	c.Next()
}
