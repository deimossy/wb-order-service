package config

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	ServerPort string

	ServerReadTimeout       time.Duration
	ServerWriteTimeout      time.Duration
	ServerIdleTimeout       time.Duration
	ServerReadHeaderTimeout time.Duration

	PgDSN             string
	PgMaxOpenConns    int
	PgMaxIdleConns    int
	PgConnMaxLifetime time.Duration

	KafkaBrokers []string
	KafkaTopic   string

	CacheCapacity int
	CacheWarmSize int
	CacheTTL      time.Duration
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig(logg *zap.Logger) *Config {
	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			logg.Info("Warning: .env file not found, using environment variables")
		}

		cfg = &Config{
			ServerPort: getEnv("SERVER_PORT", "8080"),

			ServerReadTimeout:       getEnvAsDuration("SERVER_HTTP_READ_TIMEOUT", 5*time.Second),
			ServerWriteTimeout:      getEnvAsDuration("SERVER_HTTP_WRITE_TIMEOUT", 10*time.Second),
			ServerIdleTimeout:       getEnvAsDuration("SERVER_HTTP_IDLE_TIMEOUT", 120*time.Second),
			ServerReadHeaderTimeout: getEnvAsDuration("SERVER_HTTP_READ_HEADER_TIMEOUT", 5*time.Second),

			PgDSN:             getEnv("PG_DSN", "postgresql://app:app@postgres:5432/app?sslmode=disable"),
			PgMaxOpenConns:    getEnvAsInt("PG_MAX_OPEN_CONNS", 10),
			PgMaxIdleConns:    getEnvAsInt("PG_MAX_IDLE_CONNS", 5),
			PgConnMaxLifetime: getEnvAsDuration("PG_CONN_MAX_LIFETIME", time.Hour),

			KafkaBrokers: parseCSV(getEnv("KAFKA_BROKERS", "kafka:9092")),
			KafkaTopic:   getEnv("KAFKA_TOPIC", "get_orders"),

			CacheCapacity: getEnvAsInt("CACHE_CAPACITY", 1000),
			CacheWarmSize: getEnvAsInt("CACHE_WARM_SIZE", 50),
			CacheTTL:      getEnvAsDuration("CACHE_TTL", 5*time.Minute),
		}

		logg.Info("config loaded",
			zap.String("server_port", cfg.ServerPort),
			zap.String("kafka_topic", cfg.KafkaTopic),
		)
	})

	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}

func getEnvAsInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}

	return def
}

func getEnvAsDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}

	return def
}

func parseCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	return parts
}
