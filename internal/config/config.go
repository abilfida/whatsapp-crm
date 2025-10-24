package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT
	JWTSecret    string
	JWTExpiresIn string

	// WhatsApp API
	WhatsAppAPIURL          string
	WhatsAppAPIToken        string
	WhatsAppWebhookVerifyToken string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string

	// Server
	Port string
	Env  string

	// CORS
	CORSOrigins string

	// File Upload
	MaxFileSize string
	UploadPath  string
}

func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "whatsapp_crm"),

		// JWT
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),

		// WhatsApp API
		WhatsAppAPIURL:          getEnv("WHATSAPP_API_URL", ""),
		WhatsAppAPIToken:        getEnv("WHATSAPP_API_TOKEN", ""),
		WhatsAppWebhookVerifyToken: getEnv("WHATSAPP_WEBHOOK_VERIFY_TOKEN", ""),

		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnv("REDIS_DB", "0"),

		// Server
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		// CORS
		CORSOrigins: getEnv("CORS_ORIGINS", "*"),

		// File Upload
		MaxFileSize: getEnv("MAX_FILE_SIZE", "10485760"),
		UploadPath:  getEnv("UPLOAD_PATH", "./uploads"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}