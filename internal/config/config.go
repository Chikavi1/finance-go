package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	MinIO    MinIOConfig
	Asynq    AsynqConfig
	App      AppConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	URL            string
	Host           string
	Port           string
	User           string
	Password       string
	Name           string
	SSLMode        string
	MaxConns       int
	MinConns       int
}

type RedisConfig struct {
	URL      string
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret             string
	AccessExpiration   time.Duration
	RefreshExpiration  time.Duration
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type AsynqConfig struct {
	RedisURL string
}

type AppConfig struct {
	Name             string
	Env              string
	LogLevel         string
	CORSAllowedOrigins string
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("..")

	v.AutomaticEnv()

	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("SERVER_HOST", "0.0.0.0")
	v.SetDefault("DATABASE_URL", "postgres://finances:finances@localhost:5432/finances?sslmode=disable")
	v.SetDefault("DATABASE_HOST", "localhost")
	v.SetDefault("DATABASE_PORT", "5432")
	v.SetDefault("DATABASE_USER", "finances")
	v.SetDefault("DATABASE_PASSWORD", "finances")
	v.SetDefault("DATABASE_NAME", "finances")
	v.SetDefault("DATABASE_SSLMODE", "disable")
	v.SetDefault("DATABASE_MAX_CONNS", 25)
	v.SetDefault("DATABASE_MIN_CONNS", 5)
	v.SetDefault("REDIS_URL", "redis://localhost:6379/0")
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("JWT_SECRET", "change-this-to-a-secure-random-string-at-least-32-chars")
	v.SetDefault("JWT_ACCESS_EXPIRATION", "15m")
	v.SetDefault("JWT_REFRESH_EXPIRATION", "168h")
	v.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	v.SetDefault("MINIO_ACCESS_KEY", "minio")
	v.SetDefault("MINIO_SECRET_KEY", "minio123")
	v.SetDefault("MINIO_BUCKET", "finances")
	v.SetDefault("MINIO_USE_SSL", false)
	v.SetDefault("ASYNC_REDIS_URL", "redis://localhost:6379/1")
	v.SetDefault("APP_NAME", "Finances API")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("LOG_LEVEL", "debug")
	v.SetDefault("CORS_ALLOWED_ORIGINS", "*")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	accessExp, err := time.ParseDuration(v.GetString("JWT_ACCESS_EXPIRATION"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_EXPIRATION: %w", err)
	}

	refreshExp, err := time.ParseDuration(v.GetString("JWT_REFRESH_EXPIRATION"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_EXPIRATION: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: v.GetString("SERVER_PORT"),
			Host: v.GetString("SERVER_HOST"),
		},
		Database: DatabaseConfig{
			URL:      v.GetString("DATABASE_URL"),
			Host:     v.GetString("DATABASE_HOST"),
			Port:     v.GetString("DATABASE_PORT"),
			User:     v.GetString("DATABASE_USER"),
			Password: v.GetString("DATABASE_PASSWORD"),
			Name:     v.GetString("DATABASE_NAME"),
			SSLMode:  v.GetString("DATABASE_SSLMODE"),
			MaxConns: v.GetInt("DATABASE_MAX_CONNS"),
			MinConns: v.GetInt("DATABASE_MIN_CONNS"),
		},
		Redis: RedisConfig{
			URL:      v.GetString("REDIS_URL"),
			Host:     v.GetString("REDIS_HOST"),
			Port:     v.GetString("REDIS_PORT"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:            v.GetString("JWT_SECRET"),
			AccessExpiration:  accessExp,
			RefreshExpiration: refreshExp,
		},
		MinIO: MinIOConfig{
			Endpoint:  v.GetString("MINIO_ENDPOINT"),
			AccessKey: v.GetString("MINIO_ACCESS_KEY"),
			SecretKey: v.GetString("MINIO_SECRET_KEY"),
			Bucket:    v.GetString("MINIO_BUCKET"),
			UseSSL:    v.GetBool("MINIO_USE_SSL"),
		},
		Asynq: AsynqConfig{
			RedisURL: v.GetString("ASYNC_REDIS_URL"),
		},
		App: AppConfig{
			Name:               v.GetString("APP_NAME"),
			Env:                v.GetString("APP_ENV"),
			LogLevel:           v.GetString("LOG_LEVEL"),
			CORSAllowedOrigins: v.GetString("CORS_ALLOWED_ORIGINS"),
		},
	}

	return cfg, nil
}
