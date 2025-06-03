package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Log            *LogConfig            `mapstructure:"log"`
	Db             *DBConfig             `mapstructure:"db"`
	Rd             *RedisConfig          `mapstructure:"rd"`
	LibraryService *LibraryServiceConfig `mapstructure:"library_service"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (c *DBConfig) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

type RedisConfig struct {
	Addr            string `mapstructure:"addr"`
	Password        string `mapstructure:"password"`
	DB              int    `mapstructure:"db"`
	PoolSize        int    `mapstructure:"pool_size"`
	CacheTTLSeconds int    `mapstructure:"cache_ttl_seconds"`
}

type LibraryServiceConfig struct {
	Port int `mapstructure:"port"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType(strings.TrimPrefix(path[strings.LastIndex(path, ".")+1:], ".")) // yml, json, etc.
	viper.AutomaticEnv()                                                                // Read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
