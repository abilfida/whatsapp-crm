package main

import (
	"log"
	"os"
	"whatsapp-crm/internal/models"
	"whatsapp-crm/pkg/database"
)

func main() {
	// use environment or defaults from .env.example
	cfg := struct{ DBHost, DBPort, DBUser, DBPassword, DBName string }{
		DBHost: os.Getenv("DB_HOST"), DBPort: os.Getenv("DB_PORT"), DBUser: os.Getenv("DB_USER"), DBPassword: os.Getenv("DB_PASSWORD"), DBName: os.Getenv("DB_NAME"),
	}
	// quick connect via existing helper
	db, err := database.Connect(&struct { // minimal adapter
		DBHost, DBPort, DBUser, DBPassword, DBName string
		JWTSecret, JWTExpiresIn string
		WhatsAppAPIURL, WhatsAppAPIToken, WhatsAppWebhookVerifyToken string
		RedisHost, RedisPort, RedisPassword, RedisDB string
		Port, Env, CORSOrigins, MaxFileSize, UploadPath string
	}{DBHost: cfg.DBHost, DBPort: cfg.DBPort, DBUser: cfg.DBUser, DBPassword: cfg.DBPassword, DBName: cfg.DBName})
	if err != nil { log.Fatal(err) }

	// seed admin if not exists
	var count int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count == 0 {
		// default admin: admin@example.com / Admin123!
		pass := "$2a$10$k1g7Jm3n7NQj9w6RrQx2eOq4kOe5d1uP3bY0Jv1gW7X2HfYfCk0B6" // bcrypt for Admin123!
		admin := models.User{Name: "Administrator", Email: "admin@example.com", Password: pass, Role: "admin", Status: "active"}
		if err := db.Create(&admin).Error; err != nil { log.Fatal(err) }
		log.Println("Seeded default admin: admin@example.com / Admin123!")
	} else {
		log.Println("Admin already exists; skipping seed")
	}
}
