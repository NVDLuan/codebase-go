// Package config đọc cấu hình từ biến môi trường (ConfigMap + Secret trên K8s).
package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr        string
	AppEnv          string
	LogLevel        string
	ShutdownTimeout time.Duration

	// Rỗng → dùng adapter local (in-memory / log), phù hợp chạy local hoặc dev chưa có hạ tầng.
	DatabaseURL  string   // Secret
	StripeKey    string   // Secret
	KafkaBrokers []string // ConfigMap, phân tách bằng dấu phẩy
	KafkaTopic   string
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:        env("HTTP_ADDR", ":8080"),
		AppEnv:          env("APP_ENV", "local"),
		LogLevel:        env("LOG_LEVEL", "info"),
		ShutdownTimeout: 20 * time.Second,
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		StripeKey:       os.Getenv("STRIPE_API_KEY"),
		KafkaTopic:      env("KAFKA_TOPIC", "orders"),
	}
	if b := os.Getenv("KAFKA_BROKERS"); b != "" {
		cfg.KafkaBrokers = strings.Split(b, ",")
	}
	if v := os.Getenv("SHUTDOWN_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT: %w", err)
		}
		cfg.ShutdownTimeout = d
	}
	return cfg, nil
}

func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
