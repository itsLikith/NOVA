package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(
	userID string,
	role string,
	issuer string,
	secret string,
	expirationHours int,
) (string, error) {
	now := time.Now()

	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  userID,
			Issuer:   issuer,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(
				time.Duration(expirationHours) * time.Hour,
			)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(secret))
}

func ParseToken(
	tokenString string,
	secret string,
	issuer string,
) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(secret), nil
		},
		jwt.WithIssuer(issuer),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func ExtractBearerToken(header string) (string, error) {
	parts := strings.Fields(header)

	if len(parts) != 2 {
		return "", errors.New("invalid authorization header")
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization scheme")
	}

	return parts[1], nil
}
