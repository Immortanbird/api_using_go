package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenClaims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

type TokenDetails struct {
	Token     string
	JTI       string
	ExpiresAt time.Time
}

func GenerateToken(userID string, lifespan int, jwtSecretKey string) (TokenDetails, error) {
	details := TokenDetails{}

	jti := uuid.NewString()
	iat := time.Now()
	exp := iat.Add(time.Duration(lifespan) * time.Second)

	claims := &TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			// Set the expiration time for the token.
			ExpiresAt: jwt.NewNumericDate(exp),

			// Set the time the token was issued.
			IssuedAt: jwt.NewNumericDate(iat),

			// Set the issuer of the token (your application's name or domain).
			Issuer: "api_using_go",

			// Set the subject of the token (often the user ID).
			Subject: userID,

			// Set the JWT ID.
			ID: jti,
		},
	}

	// Create a new token object, specifying the signing method (HS256) and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key to get the complete, signed token string.
	// This final string is what the client will use.
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		// If there's an error during signing, we return it.
		return details, err
	}

	details.Token = tokenString
	details.JTI = jti
	details.ExpiresAt = exp

	// Return the signed token string.
	return details, nil
}

func ParseToken(tokenString string, jwtSecretKey string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
