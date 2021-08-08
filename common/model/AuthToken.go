package model

import (
	"github.com/dgrijalva/jwt-go"
)

type AuthToken struct {
	*jwt.StandardClaims
	UserName string `json:"user_name"`
	UserID   string `json:"userID"`
	Ttl      int64  `json:"ttl"`
}
