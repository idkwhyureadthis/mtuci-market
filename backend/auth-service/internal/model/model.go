package model

import "github.com/golang-jwt/jwt/v5"

type Tokens struct {
	Access         string `json:"access"`
	Refresh        string `json:"refresh"`
	CryptedRefresh string `json:"crypted_refresh"`
}

type Token struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

type VerificationResponse struct {
	Id int `json:"id"`
}
