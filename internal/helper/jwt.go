package helper

import (
	"errors"
	"gogroceries/config"
	"gogroceries/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTInterface interface {
	GenerateToken(claims *domain.JWTClaims) (string, error)
	ValidateToken(tokenString string) (*domain.JWTClaims, error)
	ExtractJWTUser(c *gin.Context) (*domain.JWTClaims, error)
}

type jwtHelper struct {
	secretKey []byte
}

func NewJWTHelper(cfg config.Config) JWTInterface {
	return &jwtHelper{
		secretKey: []byte(cfg.JWTSecret),
	}
}

func (j *jwtHelper) GenerateToken(claims *domain.JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (j *jwtHelper) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Token is not valid")
}

func (j *jwtHelper) ExtractJWTUser(c *gin.Context) (*domain.JWTClaims, error) {
	userClaims, exists := c.Get("user_claims")
	if !exists {
		return nil, errors.New("User claims not found in context")
	}

	claims, ok := userClaims.(*domain.JWTClaims)
	if !ok {
		return nil, errors.New("Failed to cast user claims")
	}

	return claims, nil
}