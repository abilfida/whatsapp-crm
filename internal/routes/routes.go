package routes

import (
	"whatsapp-crm/internal/config"
	"whatsapp-crm/internal/controllers"
	"whatsapp-crm/internal/middlewares"
	"whatsapp-crm/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, rdb *redis.Client, cfg *config.Config) {
	api := app.Group("/api/v1")

	// Middlewares
	authMw := middlewares.NewAuthMiddleware(db)

	// Services
	customerSvc := services.NewCustomerService(db)
	conversationSvc := services.NewConversationService(db)
	messageSvc := services.NewMessageService(db, cfg)

	// Controllers
	authCtl := controllers.NewAuthController(db)
	userCtl := controllers.NewUserController(db)
	customerCtl := controllers.NewCustomerController(db, customerSvc)
	conversationCtl := controllers.NewConversationController(db, conversationSvc, messageSvc)
	messageCtl := controllers.NewMessageController(db, messageSvc)
	webhookCtl := controllers.NewWebhookController(db, cfg, messageSvc, customerSvc, conversationSvc)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", authCtl.Login)
	auth.Post("/register", authCtl.Register)
	auth.Get("/profile", authMw.RequireAuth, authCtl.GetProfile)
	auth.Post("/change-password", authMw.RequireAuth, authCtl.ChangePassword)

	// Users (admin only)
	users := api.Group("/users", authMw.RequireAuth, authMw.RequireRole("admin"))
	users.Get("/", userCtl.GetUsers)
	users.Get("/:id", userCtl.GetUser)
	users.Post("/", userCtl.CreateUser)
	users.Put("/:id", userCtl.UpdateUser)
	users.Delete("/:id", userCtl.DeleteUser)

	// Customers
	customers := api.Group("/customers", authMw.RequireAuth)
	customers.Get("/", customerCtl.List)
	customers.Post("/", customerCtl.Create)
	customers.Get("/:id", customerCtl.Detail)
	customers.Put("/:id", customerCtl.Update)
	customers.Delete("/:id", customerCtl.Delete)

	// Conversations
	convs := api.Group("/conversations", authMw.RequireAuth)
	convs.Get("/", conversationCtl.List)
	convs.Post("/", conversationCtl.Create)
	convs.Get("/:id", conversationCtl.Detail)
	convs.Put("/:id/assign", conversationCtl.Assign)
	convs.Put("/:id/status", conversationCtl.UpdateStatus)
	convs.Put("/:id/priority", conversationCtl.UpdatePriority)
	convs.Put("/:id/notes", conversationCtl.UpdateNotes)

	// Messages
	msgs := api.Group("/messages", authMw.RequireAuth)
	msgs.Get("/conversation/:id", messageCtl.ListByConversation)
	msgs.Post("/conversation/:id/text", messageCtl.SendText)
	msgs.Post("/conversation/:id/media", messageCtl.SendMedia)
	msgs.Post("/conversation/:id/template", messageCtl.SendTemplate)

	// Webhook (note: user asked to hold execution usage; routes exist but you can disable mount if needed)
	webhook := api.Group("/webhook")
	webhook.Get("/whatsapp", webhookCtl.VerifyWebhook)
	webhook.Post("/whatsapp", webhookCtl.HandleWebhook)
}
