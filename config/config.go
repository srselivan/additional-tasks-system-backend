package config

import (
	"sync"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPServer HTTPServer
	Logger     Logger
	JWT        JWT
}

type Logger struct {
	Level string `env:"LOGGER_LEVEL" envDefault:"info"`
}

type HTTPServer struct {
	Addr string `env:"HTTP_SERVER_ADDR" envDefault:"0.0.0.0:11225"`
}

type JWT struct {
	JWTRefreshTokenExpTime time.Duration `env:"JWT_REFRESH_EXPIRATION_TIME,required"`
	JWTAccessTokenExpTime  time.Duration `env:"JWT_ACCESS_EXPIRATION_TIME,required"`
	JWTRefreshSecretKey    string        `env:"JWT_REFRESH_SECRET_KEY,required"`
	JWTAccessSecretKey     string        `env:"JWT_ACCESS_SECRET_KEY,required"`
}

var (
	config Config
	once   sync.Once
)

func New() *Config {
	once.Do(func() {
		_ = godotenv.Load()
		if err := env.Parse(&config); err != nil {
			panic(err)
		}
	})
	return &config
}
