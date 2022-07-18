package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type JWTService struct {
	secretKey string
	issuer    string
}

type jwtCustomClaim struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func NewJWTService() *JWTService {
	return &JWTService{
		issuer:    "ChoTot",
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	if err := godotenv.Load(".env"); err != nil {
		panic("Error loading .env file")
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey != "" {
		secretKey = "ChoTotKey"
	}
	return secretKey
}

func (j *JWTService) GenerateToken(userId string) string {
	claims := &jwtCustomClaim{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (j *JWTService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}

func (j *JWTService) ParseToken(authHeader string) (int, error) {
	token, err := j.ValidateToken(authHeader)
	if err != nil {
		return -1, err
	}
	claims := token.Claims.(jwt.MapClaims)
	userId, _ := strconv.Atoi(claims["user_id"].(string))
	return userId, nil
}
