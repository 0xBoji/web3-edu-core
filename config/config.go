package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/ini.v1"
)

type Server struct {
	RunMode      string
	Host         string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	BuildName    string
}

type Database struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	TimeZone string
}

type App struct {
	Name                   string
	JWTSecret              string
	TokenExpireTime        int
	RefreshTokenExpireTime int
}

type Redis struct {
	Host     string
	Port     int
	Password string
	DB       int
}

var (
	ServerSetting   = &Server{}
	DatabaseSetting = &Database{}
	AppSetting      = &App{}
	RedisSetting    = &Redis{}
)

// Setup initializes the configuration instance
func Setup() {
	// Load configuration from file
	cfg, err := ini.Load("config/app.ini")
	if err != nil {
		log.Printf("Warning: Failed to parse 'config/app.ini': %v", err)
		log.Println("Falling back to environment variables")
	}

	// Map configuration from file
	if cfg != nil {
		mapTo(cfg, "server", ServerSetting)
		mapTo(cfg, "database", DatabaseSetting)
		mapTo(cfg, "app", AppSetting)
		mapTo(cfg, "redis", RedisSetting)
	}

	// Override with environment variables if they exist
	loadEnvVariables()

	// Set timeouts
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
}

// mapTo maps section to struct
func mapTo(cfg *ini.File, section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Printf("Warning: Failed to map section '%s': %v", section, err)
	}
}

// loadEnvVariables loads configuration from environment variables
func loadEnvVariables() {
	// Server settings
	if env := os.Getenv("SERVER_RUN_MODE"); env != "" {
		ServerSetting.RunMode = env
	}
	if env := os.Getenv("SERVER_HOST"); env != "" {
		ServerSetting.Host = env
	}
	if env := os.Getenv("SERVER_HTTP_PORT"); env != "" {
		if port, err := strconv.Atoi(env); err == nil {
			ServerSetting.HttpPort = port
		}
	}
	if env := os.Getenv("SERVER_READ_TIMEOUT"); env != "" {
		if timeout, err := strconv.Atoi(env); err == nil {
			ServerSetting.ReadTimeout = time.Duration(timeout)
		}
	}
	if env := os.Getenv("SERVER_WRITE_TIMEOUT"); env != "" {
		if timeout, err := strconv.Atoi(env); err == nil {
			ServerSetting.WriteTimeout = time.Duration(timeout)
		}
	}
	if env := os.Getenv("SERVER_BUILD_NAME"); env != "" {
		ServerSetting.BuildName = env
	}

	// Database settings
	if env := os.Getenv("POSTGRES_HOST"); env != "" {
		DatabaseSetting.Host = env
	}
	if env := os.Getenv("POSTGRES_PORT"); env != "" {
		if port, err := strconv.Atoi(env); err == nil {
			DatabaseSetting.Port = port
		}
	}
	if env := os.Getenv("POSTGRES_USER"); env != "" {
		DatabaseSetting.User = env
	}
	if env := os.Getenv("POSTGRES_PASSWORD"); env != "" {
		DatabaseSetting.Password = env
	}
	if env := os.Getenv("POSTGRES_DB"); env != "" {
		DatabaseSetting.Name = env
	}
	if env := os.Getenv("POSTGRES_SSL_MODE"); env != "" {
		DatabaseSetting.SSLMode = env
	}
	if env := os.Getenv("POSTGRES_TIMEZONE"); env != "" {
		DatabaseSetting.TimeZone = env
	}

	// App settings
	if env := os.Getenv("APP_NAME"); env != "" {
		AppSetting.Name = env
	}
	if env := os.Getenv("APP_JWT_SECRET"); env != "" {
		AppSetting.JWTSecret = env
	}
	if env := os.Getenv("APP_TOKEN_EXPIRE_TIME"); env != "" {
		if time, err := strconv.Atoi(env); err == nil {
			AppSetting.TokenExpireTime = time
		}
	}
	if env := os.Getenv("APP_REFRESH_TOKEN_EXPIRE_TIME"); env != "" {
		if time, err := strconv.Atoi(env); err == nil {
			AppSetting.RefreshTokenExpireTime = time
		}
	}

	// Redis settings
	if env := os.Getenv("REDIS_HOST"); env != "" {
		RedisSetting.Host = env
	}
	if env := os.Getenv("REDIS_PORT"); env != "" {
		if port, err := strconv.Atoi(env); err == nil {
			RedisSetting.Port = port
		}
	}
	if env := os.Getenv("REDIS_PASSWORD"); env != "" {
		RedisSetting.Password = env
	}
	if env := os.Getenv("REDIS_DB"); env != "" {
		if db, err := strconv.Atoi(env); err == nil {
			RedisSetting.DB = db
		}
	}
}
