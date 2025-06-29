package handler

import (
	"net/http"

	"github.com/Immortanbird/api_using_go/internal/models"
	"github.com/Immortanbird/api_using_go/internal/repository/crud"
	"github.com/Immortanbird/api_using_go/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// POST	/auth/register
func (h *Handler) Register(c *gin.Context) {
	payload := models.RegisterRequest{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		zap.L().Warn(
			"Failed to bind JSON request for registration",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and Password are required."})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Warn("Failed to hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user."})
		return
	}

	user := models.Users{
		Email:    payload.Email,
		Username: payload.Username,
		Password: string(hashed),
	}

	err = crud.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user."})
		return
	}

	zap.L().Info(
		"Registration request received",
		zap.String("client_ip", c.ClientIP()),
		zap.String("path", c.Request.URL.Path),
		zap.String("email", payload.Email),
		zap.String("username", payload.Username),
	)

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// POST	/auth/login
func (h *Handler) Login(c *gin.Context) {
	payload := models.LoginRequest{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		zap.L().Warn(
			"Failed to bind JSON request for login",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and Password are required."})
		return
	}

	// Fetch user from database
	var user *models.Users
	if payload.Email != "" {
		var err error
		user, err = crud.FindTheUser(models.Users{Email: payload.Email})
		if err != nil {
			zap.L().Warn(
				"Failed to find user by email",
				zap.Error(err),
				zap.String("client_ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
				zap.String("email", payload.Email),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password."})
			return
		}
	} else if payload.Username != "" {
		var err error
		user, err = crud.FindTheUser(models.Users{Username: payload.Username})
		if err != nil {
			zap.L().Warn(
				"Failed to find user by username",
				zap.Error(err),
				zap.String("client_ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
				zap.String("email", payload.Email),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password."})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Username is required."})
		return
	}

	// Wrong password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		zap.L().Warn(
			"Password mismatch during login",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("email", (*user).ID.String()),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password."})
		return
	}

	// Generate access and refresh tokens
	accessToken, err := utils.GenerateToken(user.ID.String(), h.Config.JWT.ExpAccess, h.Config.JWT.SecretKey)
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
	refreshToken, err := utils.GenerateToken(user.ID.String(), h.Config.JWT.ExpRefresh, h.Config.JWT.SecretKey)
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
		UserID:    user.ID,
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
			zap.String("user_id", user.ID.String()),
			zap.String("refresh_token_jti", refreshToken.JTI),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token."})
		return
	}

	c.SetCookie("access_token", accessToken.Token, h.Config.JWT.ExpAccess, "/", "", true, true)
	c.SetCookie("refresh_token", refreshToken.Token, h.Config.JWT.ExpRefresh, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}

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
