package middlewares

import (
	"net/http"

	"github.com/Immortanbird/api_using_go/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

	userID, err := uuid.Parse(claims.ID)

	if err != nil {
		zap.L().Error(
			"Unable to parse user ID",
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
		)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Internal Server Error."})
		return
	}

	c.Set("userID", userID)

	c.Next()
}
