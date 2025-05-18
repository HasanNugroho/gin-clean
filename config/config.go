package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   Server   `mapstructure:",squash"`
	Database Database `mapstructure:",squash"`
	Secret   Secret   `mapstructure:",squash"`
	Redis    Redis    `mapstructure:",squash"`
	Context  Context  `mapstructure:",squash"`
	Security Security `mapstructure:",squash"`
}

type Server struct {
	Name     string `mapstructure:"APP_NAME"`
	Host     string `mapstructure:"APP_HOST"`
	Port     string `mapstructure:"APP_PORT"`
	BaseUrl  string `mapstructure:"APP_BASEURL"`
	LogLevel int    `mapstructure:"LOG_LEVEL"`
}

type Database struct {
	Host string `mapstructure:"POSTGRES_HOST"`
	Port string `mapstructure:"POSTGRES_PORT"`
	User string `mapstructure:"POSTGRES_USER"`
	Pass string `mapstructure:"POSTGRES_PASS"`
	Name string `mapstructure:"POSTGRES_DB"`
	Ssl  string `mapstructure:"POSTGRES_SSL"`
}

type Secret struct {
	Jwt                string `mapstructure:"SECRET_KEY"`
	TokenExpiry        string `mapstructure:"TOKEN_EXPIRY"`
	RefreshTokenExpiry string `mapstructure:"REFRESH_TOKEN_EXPIRY"`
}

type Redis struct {
	Host string `mapstructure:"REDIS_HOST"`
	Port string `mapstructure:"REDIS_PORT"`
	Pass string `mapstructure:"REDIS_PASS"`
}

type Context struct {
	Timeout int `mapstructure:"TIMEOUT"`
}

type Security struct {
	RateLimit         string   `mapstructure:"RATE_LIMIT"`
	AllowedOrigins    []string `mapstructure:"ALLOWED_ORIGINS"`
	TrustedPlatform   string   `mapstructure:"TRUSTED_PLATFORM"`
	ExpectedHost      string   `mapstructure:"EXPECTED_HOST"`
	XFrameOptions     string   `mapstructure:"X_FRAME_OPTIONS"`
	ContentSecurity   string   `mapstructure:"CONTENT_SECURITY_POLICY"`
	XXSSProtection    string   `mapstructure:"X_XSS_PROTECTION"`
	StrictTransport   string   `mapstructure:"STRICT_TRANSPORT_SECURITY"`
	ReferrerPolicy    string   `mapstructure:"REFERRER_POLICY"`
	XContentTypeOpts  string   `mapstructure:"X_CONTENT_TYPE_OPTIONS"`
	PermissionsPolicy string   `mapstructure:"PERMISSIONS_POLICY"`
}

func Get() (config *Config, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("no .env file found, using system environment variables: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if config.Server.LogLevel < -1 || config.Server.LogLevel > 5 {
		return nil, fmt.Errorf("LOG_LEVEL must be between -1 (trace) and 5 (panic), got %d", config.Server.LogLevel)
	}

	if _, err := time.ParseDuration(config.Secret.TokenExpiry); err != nil {
		return nil, fmt.Errorf("invalid TOKEN_EXPIRY: %w", err)
	}

	if _, err := time.ParseDuration(config.Secret.RefreshTokenExpiry); err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRY: %w", err)
	}

	timeoutSeconds := viper.GetInt("TIMEOUT")
	if timeoutSeconds <= 0 {
		timeoutSeconds = 3600
	}

	config.Context.Timeout = int(time.Duration(timeoutSeconds) * time.Second)
	config.Security.AllowedOrigins = strings.Split(viper.GetString("ALLOWED_ORIGINS"), ",")

	return
}
