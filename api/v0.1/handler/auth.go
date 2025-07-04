package handler

import (
	"net/http"

	"github.com/Immortanbird/api_using_go/internal/models"
	"github.com/Immortanbird/api_using_go/internal/repository/crud"
	"github.com/Immortanbird/api_using_go/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// POST	/auth/refresh
func (h *Handler) Refresh(c *gin.Context) {
	tokenString, err := c.Cookie("refresh_token")
	if err != nil {
		zap.L().Warn(
			"No refresh token found in the cookie.",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
		)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Missing token."})
	}

	claims, err := utils.ParseToken(tokenString, h.Config.JWT.SecretKey)
	if err != nil {
		zap.L().Warn(
			"Failed to parse token",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
		)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid token."})
		return
	}

	token, err := crud.FindRefreshToken(models.RefreshToken{TokenJTI: claims.ID})
	if err != nil || token.IsRevoked {
		zap.L().Warn(
			"Failed to find refresh token in the database",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("user_id", claims.UserID),
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid refresh token."})
		return
	}

	// Generate new access and refresh tokens
	accessToken, err := utils.GenerateToken(token.UserID.String(), h.Config.JWT.ExpAccess, h.Config.JWT.SecretKey)
	if err != nil {
		zap.L().Warn(
			"Failed to generate token",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token."})
		return
	}
	refreshToken, err := utils.GenerateToken(token.UserID.String(), h.Config.JWT.ExpRefresh, h.Config.JWT.SecretKey)
	if err != nil {
		zap.L().Warn(
			"Failed to generate token",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token."})
		return
	}

	// Save the refresh token in the database
	err = crud.CreateRefreshToken(&models.RefreshToken{
		UserID:    token.UserID,
		TokenJTI:  refreshToken.JTI,
		IsRevoked: false,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		ExpiresAt: refreshToken.ExpiresAt,
	})
	if err != nil {
		zap.L().Warn(
			"Failed to create refresh token in the database",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("user_id", token.UserID.String()),
			zap.String("refresh_token_jti", refreshToken.JTI),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token."})
		return
	}

	c.SetCookie("access_token", accessToken.Token, h.Config.JWT.ExpAccess, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken.Token, h.Config.JWT.ExpRefresh, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed."})
}
