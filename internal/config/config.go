package config

import (
	"log"
	"os"
	"strconv"

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
	WhatsAppAPIURL             string
	WhatsAppAPIToken           string
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
	MaxFileSize int64
	UploadPath  string
	PublicBaseURL string

	// Storage
	StorageDriver               string
	StorageBasePath            string
	StorageSignedURLExpSeconds int64

	// AWS
	AWSRegion              string
	AWSS3Bucket            string
	AWSS3Prefix            string
	AWSAccessKeyID         string
	AWSSecretAccessKey     string
	AWSS3SignedURLExpSeconds int64

	// GCS
	GCSBucket              string
	GCSPrefix              string
	GCredentialsPath       string
	GCSSignedURLExpSeconds int64

	// Allowed types
	AllowedImageTypes    []string
	AllowedDocumentTypes []string
	AllowedAudioTypes    []string
	AllowedVideoTypes    []string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	maxSize := int64(parseInt("UPLOAD_MAX_SIZE", 16777216))
	signedExp := int64(parseInt("STORAGE_SIGNED_URL_EXP_SECONDS", 86400))
	s3Exp := int64(parseInt("AWS_S3_SIGNED_URL_EXP_SECONDS", 86400))
	gcsExp := int64(parseInt("GCS_SIGNED_URL_EXP_SECONDS", 86400))

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "whatsapp_crm"),

		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),

		WhatsAppAPIURL:             getEnv("WHATSAPP_API_URL", ""),
		WhatsAppAPIToken:           getEnv("WHATSAPP_API_TOKEN", ""),
		WhatsAppWebhookVerifyToken: getEnv("WHATSAPP_WEBHOOK_VERIFY_TOKEN", ""),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnv("REDIS_DB", "0"),

		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		CORSOrigins: getEnv("CORS_ORIGINS", "*"),

		MaxFileSize: maxSize,
		UploadPath:  getEnv("UPLOAD_PATH", "./uploads"),
		PublicBaseURL: getEnv("PUBLIC_BASE_URL", ""),

		StorageDriver:               getEnv("STORAGE_DRIVER", "local"),
		StorageBasePath:            getEnv("STORAGE_BASE_PATH", "./uploads/media"),
		StorageSignedURLExpSeconds: signedExp,

		AWSRegion:              getEnv("AWS_REGION", ""),
		AWSS3Bucket:            getEnv("AWS_S3_BUCKET", ""),
		AWSS3Prefix:            getEnv("AWS_S3_PREFIX", "whatsapp-crm"),
		AWSAccessKeyID:         getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey:     getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSS3SignedURLExpSeconds: s3Exp,

		GCSBucket:              getEnv("GCS_BUCKET", ""),
		GCSPrefix:              getEnv("GCS_PREFIX", "whatsapp-crm"),
		GCredentialsPath:       getEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),
		GCSSignedURLExpSeconds: gcsExp,

		AllowedImageTypes:    splitCSV(getEnv("ALLOWED_IMAGE_TYPES", "jpg,jpeg,png,gif,webp")),
		AllowedDocumentTypes: splitCSV(getEnv("ALLOWED_DOCUMENT_TYPES", "pdf,doc,docx,xls,xlsx,ppt,pptx")),
		AllowedAudioTypes:    splitCSV(getEnv("ALLOWED_AUDIO_TYPES", "mp3,ogg,m4a,wav,aac")),
		AllowedVideoTypes:    splitCSV(getEnv("ALLOWED_VIDEO_TYPES", "mp4,3gp,mov,avi,mkv")),
	}
}

func getEnv(key, def string) string { if v := os.Getenv(key); v != "" { return v }; return def }

func parseInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil { return n }
	}
	return def
}

func splitCSV(s string) []string {
	if s == "" { return nil }
	s := s
	var out []string
	buf := ""
	for i := 0; i < len(s); i++ {
		if s[i] == ',' { if buf != "" { out = append(out, buf); buf = "" }; continue }
		if s[i] == ' ' { continue }
		buf += string(s[i])
	}
	if buf != "" { out = append(out, buf) }
	return out
}
