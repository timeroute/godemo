package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port string
	Mode string // debug, release, test
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type JWTConfig struct {
	Secret     string
	ExpireHour int
}

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int // 预检请求缓存时间（秒）
}

type LoggingConfig struct {
	SaveToDB        bool
	LogToConsole    bool
	LogRequestBody  bool
	LogResponseBody bool
	MaxBodySize     int
	GeoIPDBPath     string
}

var AppConfig *Config

// Load 加载配置
// 优先级: 环境变量 > .env 文件 > 默认值
func Load() {
	// 尝试加载 .env 文件（如果存在）
	// 不存在也不报错，因为可以直接使用环境变量
	_ = godotenv.Load()

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"), // debug, release, test
		},
		Database: DatabaseConfig{
			Driver: getEnv("DB_DRIVER", "sqlite"),
			DSN:    getEnv("DB_DSN", "godemo.db"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpireHour: getEnvAsInt("JWT_EXPIRE_HOUR", 24),
		},
		CORS: CORSConfig{
			AllowOrigins:     getEnvAsSlice("CORS_ALLOW_ORIGINS", []string{"*"}),
			AllowMethods:     getEnvAsSlice("CORS_ALLOW_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowHeaders:     getEnvAsSlice("CORS_ALLOW_HEADERS", []string{"Origin", "Content-Type", "Authorization"}),
			ExposeHeaders:    getEnvAsSlice("CORS_EXPOSE_HEADERS", []string{"Content-Length"}),
			AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           getEnvAsInt("CORS_MAX_AGE", 43200), // 12小时
		},
		Logging: LoggingConfig{
			SaveToDB:        getEnvAsBool("LOG_SAVE_TO_DB", true),
			LogToConsole:    getEnvAsBool("LOG_TO_CONSOLE", true),
			LogRequestBody:  getEnvAsBool("LOG_REQUEST_BODY", false),
			LogResponseBody: getEnvAsBool("LOG_RESPONSE_BODY", false),
			MaxBodySize:     getEnvAsInt("LOG_MAX_BODY_SIZE", 1024),
			GeoIPDBPath:     getEnv("GEOIP_DB_PATH", ""),
		},
	}

	// 验证必要的配置
	if AppConfig.JWT.Secret == "your-secret-key-change-in-production" {
		log.Println("⚠️  警告: 使用默认 JWT 密钥，生产环境请设置 JWT_SECRET 环境变量")
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取整数类型的环境变量
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsBool 获取布尔类型的环境变量
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsSlice 获取字符串切片类型的环境变量（逗号分隔）
func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	// 按逗号分隔并去除空格
	var result []string
	for _, v := range splitAndTrim(valueStr, ",") {
		if v != "" {
			result = append(result, v)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}
	return result
}

// splitAndTrim 分隔字符串并去除每个元素的空格
func splitAndTrim(s, sep string) []string {
	parts := []string{}
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		parts = append(parts, trimmed)
	}
	return parts
}

// splitString 简单的字符串分隔函数
func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	var current string

	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, current)
			current = ""
			i += len(sep) - 1
		} else {
			current += string(s[i])
		}
	}
	result = append(result, current)

	return result
}

// trimSpace 去除字符串首尾空格
func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}
