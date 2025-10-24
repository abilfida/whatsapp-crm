package controllers

import (
	"encoding/json"
	"log"
	"time"
	"whatsapp-crm/internal/config"
	"whatsapp-crm/internal/models"
	"whatsapp-crm/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type WebhookController struct {
	db               *gorm.DB
	cfg              *config.Config
	messageService   *services.MessageService
	customerService  *services.CustomerService
	conversationService *services.ConversationService
}

type WebhookMessage struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Text      *struct {
		Body string `json:"body"`
	} `json:"text,omitempty"`
	Image *struct {
		URL     string `json:"url"`
		Caption string `json:"caption"`
	} `json:"image,omitempty"`
	Document *struct {
		URL      string `json:"url"`
		Filename string `json:"filename"`
		Caption  string `json:"caption"`
	} `json:"document,omitempty"`
	Audio *struct {
		URL string `json:"url"`
	} `json:"audio,omitempty"`
	Video *struct {
		URL     string `json:"url"`
		Caption string `json:"caption"`
	} `json:"video,omitempty"`
	Location *struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Name      string  `json:"name"`
	} `json:"location,omitempty"`
	Contact *struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	} `json:"contact,omitempty"`
}

type WebhookStatus struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	To        string    `json:"to"`
}

type WebhookPresence struct {
	From      string    `json:"from"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type WebhookPayload struct {
	Type     string           `json:"type"`
	Message  *WebhookMessage  `json:"message,omitempty"`
	Status   *WebhookStatus   `json:"status,omitempty"`
	Presence *WebhookPresence `json:"presence,omitempty"`
}

func NewWebhookController(db *gorm.DB, cfg *config.Config, messageService *services.MessageService, customerService *services.CustomerService, conversationService *services.ConversationService) *WebhookController {
	return &WebhookController{
		db:                  db,
		cfg:                 cfg,
		messageService:      messageService,
		customerService:     customerService,
		conversationService: conversationService,
	}
}

// VerifyWebhook verifies webhook for WhatsApp setup
func (wc *WebhookController) VerifyWebhook(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == wc.cfg.WhatsAppWebhookVerifyToken {
		return c.SendString(challenge)
	}

	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error": "Verification failed",
	})
}

// HandleWebhook processes incoming WhatsApp webhooks
func (wc *WebhookController) HandleWebhook(c *fiber.Ctx) error {
	// Log webhook payload
	payload := string(c.Body())
	webhookLog := models.WebhookLog{
		EventType: models.WebhookEventMessage, // Will be updated based on actual type
		Status:    models.WebhookStatusReceived,
		Payload:   payload,
	}
	wc.db.Create(&webhookLog)

	// Parse webhook payload
	var webhookPayload WebhookPayload
	if err := json.Unmarshal(c.Body(), &webhookPayload); err != nil {
		log.Printf("Failed to parse webhook payload: %v", err)
		webhookLog.Status = models.WebhookStatusFailed
		webhookLog.ErrorMessage = err.Error()
		wc.db.Save(&webhookLog)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload",
		})
	}

	// Process based on webhook type
	switch webhookPayload.Type {
	case "message":
		if err := wc.handleIncomingMessage(webhookPayload.Message, &webhookLog); err != nil {
			log.Printf("Failed to process incoming message: %v", err)
			webhookLog.Status = models.WebhookStatusFailed
			webhookLog.ErrorMessage = err.Error()
		} else {
			webhookLog.Status = models.WebhookStatusProcessed
		}

	case "status":
		if err := wc.handleMessageStatus(webhookPayload.Status, &webhookLog); err != nil {
			log.Printf("Failed to process message status: %v", err)
			webhookLog.Status = models.WebhookStatusFailed
			webhookLog.ErrorMessage = err.Error()
		} else {
			webhookLog.Status = models.WebhookStatusProcessed
		}

	case "presence":
		if err := wc.handlePresenceUpdate(webhookPayload.Presence, &webhookLog); err != nil {
			log.Printf("Failed to process presence update: %v", err)
			webhookLog.Status = models.WebhookStatusFailed
			webhookLog.ErrorMessage = err.Error()
		} else {
			webhookLog.Status = models.WebhookStatusProcessed
		}

	default:
		webhookLog.Status = models.WebhookStatusIgnored
		webhookLog.ErrorMessage = "Unknown event type: " + webhookPayload.Type
	}

	// Update webhook log
	now := time.Now()
	webhookLog.ProcessedAt = &now
	wc.db.Save(&webhookLog)

	return c.JSON(fiber.Map{"status": "ok"})
}

func (wc *WebhookController) handleIncomingMessage(msg *WebhookMessage, webhookLog *models.WebhookLog) error {
	webhookLog.EventType = models.WebhookEventMessage

	// Get or create customer
	customer, err := wc.customerService.GetOrCreateCustomer(msg.From, "")
	if err != nil {
		return err
	}

	// Get or create conversation
	conversation, err := wc.conversationService.GetOrCreateConversation(customer.ID)
	if err != nil {
		return err
	}

	// Create message
	message := models.Message{
		ConversationID: conversation.ID,
		WhatsAppID:     msg.ID,
		Type:           models.MessageType(msg.Type),
		Direction:      models.MessageDirectionInbound,
		Status:         models.MessageStatusDelivered,
		SentAt:         &msg.Timestamp,
		DeliveredAt:    &msg.Timestamp,
	}

	// Set message content based on type
	switch msg.Type {
	case "text":
		if msg.Text != nil {
			message.Content = msg.Text.Body
		}
	case "image":
		if msg.Image != nil {
			message.MediaURL = msg.Image.URL
			message.Caption = msg.Image.Caption
		}
	case "document":
		if msg.Document != nil {
			message.MediaURL = msg.Document.URL
			message.FileName = msg.Document.Filename
			message.Caption = msg.Document.Caption
		}
	case "audio":
		if msg.Audio != nil {
			message.MediaURL = msg.Audio.URL
		}
	case "video":
		if msg.Video != nil {
			message.MediaURL = msg.Video.URL
			message.Caption = msg.Video.Caption
		}
	case "location":
		if msg.Location != nil {
			message.Latitude = msg.Location.Latitude
			message.Longitude = msg.Location.Longitude
			message.LocationName = msg.Location.Name
		}
	case "contact":
		if msg.Contact != nil {
			message.ContactName = msg.Contact.Name
			message.ContactPhone = msg.Contact.Phone
		}
	}

	// Save message
	if err := wc.db.Create(&message).Error; err != nil {
		return err
	}

	// Update conversation
	now := time.Now()
	conversation.LastMessageAt = &now
	conversation.UnreadCount++
	wc.db.Save(&conversation)

	// Update customer last seen
	customer.LastSeen = &now
	wc.db.Save(&customer)

	return nil
}

func (wc *WebhookController) handleMessageStatus(status *WebhookStatus, webhookLog *models.WebhookLog) error {
	webhookLog.EventType = models.WebhookEventMessageStatus

	// Update message status
	var message models.Message
	if err := wc.db.First(&message, "whatsapp_id = ?", status.ID).Error; err != nil {
		// Message not found, might be external message
		return nil
	}

	// Update status and timestamp
	message.Status = models.MessageStatus(status.Status)
	switch status.Status {
	case "delivered":
		message.DeliveredAt = &status.Timestamp
	case "read":
		message.ReadAt = &status.Timestamp
	}

	return wc.db.Save(&message).Error
}

func (wc *WebhookController) handlePresenceUpdate(presence *WebhookPresence, webhookLog *models.WebhookLog) error {
	webhookLog.EventType = models.WebhookEventPresence

	// Update customer last seen
	var customer models.Customer
	if err := wc.db.First(&customer, "whatsapp_id = ?", presence.From).Error; err != nil {
		// Customer not found, ignore
		return nil
	}

	customer.LastSeen = &presence.Timestamp
	return wc.db.Save(&customer).Error
}