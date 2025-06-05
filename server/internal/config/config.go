package config

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog"
)

// --- Общие структуры ---

type Logger struct {
	LogFilePath  string        `yaml:"logFilePath"`
	LoggingLevel zerolog.Level `yaml:"loggingLevel"` // zerolog.Level может быть строкой ("info", "debug") или числом
	DevMode      bool          `yaml:"devMode"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func (db *Database) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		db.User, db.Password, db.Host, db.Port, db.Database,
	)
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password,omitempty"` // Пароль может быть опциональным
}

// --- Конфигурации для сервиса ---

// SessionConfig содержит настройки сессии.
type SessionConfig struct {
	Secret     string        `yaml:"secret"`
	Name       string        `yaml:"name"`
	TTL        time.Duration `yaml:"ttl"` // go-yaml напрямую парсит строки типа "24h", "30m"
	CookiePath string        `yaml:"cookiePath"`
	Secure     bool          `yaml:"secure"`
	HttpOnly   bool          `yaml:"httpOnly"`
	SameSite   string        `yaml:"sameSite"` // "Lax", "Strict", "None"
}

// LibraryServiceConfig содержит настройки для сервиса библиотеки
type LibraryServiceConfig struct {
	Port          int            `yaml:"port"`
	DevMode       bool           `yaml:"devMode"`
	AllowedOrigin string         `yaml:"allowedOrigin,omitempty"`
	Session       *SessionConfig `yaml:"session,omitempty"`
}

// --- Основная структура Config ---

type Config struct {
	Log     *Logger               `yaml:"logger"`
	Db      *Database             `yaml:"database"`
	Rd      *Redis                `yaml:"redis"`
	Library *LibraryServiceConfig `yaml:"library"`
}

// Load читает YAML-конфиг, указан­ный в path, и возвращает готовую Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("can't open %s: %w", path, err)
	}

	var cfg Config
	// Строгая проверка конфигурации
	if err := yaml.UnmarshalWithOptions(data, &cfg, yaml.Strict()); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	// Проверка обязательных секций конфигурации
	if cfg.Log == nil {
		return nil, fmt.Errorf("logger configuration is missing")
	}
	if cfg.Db == nil {
		return nil, fmt.Errorf("database configuration is missing")
	}
	if cfg.Rd == nil {
		return nil, fmt.Errorf("redis configuration is missing")
	}

	// Проверка настроек библиотечного сервиса
	if cfg.Library != nil {
		if cfg.Library.Port == 0 {
			return nil, fmt.Errorf("library service port must be configured")
		}
		if cfg.Library.Session != nil {
			if cfg.Library.Session.Secret == "" {
				return nil, fmt.Errorf("library session secret is required")
			}
			if cfg.Library.Session.TTL == 0 {
				return nil, fmt.Errorf("library session TTL is required and must be a valid duration string (e.g., '24h', '30m')")
			}
		}
	} else {
		return nil, fmt.Errorf("library service configuration is missing")
	}

	return &cfg, nil
}
