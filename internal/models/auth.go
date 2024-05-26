package models

import "time"

type Credentials struct {
	Email    string
	Password string
}

type JWTPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type JWTConfig struct {
	JWTAccessExpirationTime  time.Duration
	JWTRefreshExpirationTime time.Duration
	JWTAccessSecretKey       string
	JWTRefreshSecretKey      string
}
