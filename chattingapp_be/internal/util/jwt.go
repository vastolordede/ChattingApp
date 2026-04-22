package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey       string
	accessDuration  time.Duration
	refreshDuration time.Duration
}

type UserClaims struct {
	UserID    int64  `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func NewJWTManager(secret string, accessHours int, refreshHours int) *JWTManager {
	return &JWTManager{
		secretKey:       secret,
		accessDuration:  time.Duration(accessHours) * time.Hour,
		refreshDuration: time.Duration(refreshHours) * time.Hour,
	}
}

func (j *JWTManager) GenerateAccessToken(userID int64) (string, error) {
	return j.generate(userID, "access", j.accessDuration)
}

func (j *JWTManager) GenerateRefreshToken(userID int64) (string, error) {
	return j.generate(userID, "refresh", j.refreshDuration)
}

func (j *JWTManager) generate(userID int64, tokenType string, duration time.Duration) (string, error) {
	now := time.Now()

	claims := UserClaims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTManager) VerifyAccessToken(tokenStr string) (*UserClaims, error) {
	claims, err := j.verify(tokenStr)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "access" {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}

func (j *JWTManager) VerifyRefreshToken(tokenStr string) (*UserClaims, error) {
	claims, err := j.verify(tokenStr)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}

func (j *JWTManager) Verify(tokenStr string) (*UserClaims, error) {
	return j.VerifyAccessToken(tokenStr)
}

func (j *JWTManager) verify(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(j.secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}