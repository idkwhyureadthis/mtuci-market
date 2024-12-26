package model

import "github.com/golang-jwt/jwt/v5"

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type Token struct {
	Type string `json:"type"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type VerificationResponse struct {
	Id   int    `json:"id"`
	Role string `json:"role"`
}
