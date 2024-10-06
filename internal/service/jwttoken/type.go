package jwttoken

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TokenLifeTime = time.Hour * 3

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}
