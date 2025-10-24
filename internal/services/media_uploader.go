package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
	"whatsapp-crm/internal/config"
	"whatsapp-crm/internal/models"
	"whatsapp-crm/internal/storage"
	"whatsapp-crm/pkg/whatsapp"

	"github.com/google/uuid"
)

type MediaUploader struct {
	store storage.Storage
	wa    *whatsapp.Client
	expiry time.Duration
}

func NewMediaUploader(store storage.Storage, cfg *config.Config) *MediaUploader {
	return &MediaUploader{store: store, wa: whatsapp.NewClient(cfg), expiry: 24 * time.Hour}
}

func detectType(filename, headerCT, declaredType string) (mediaType, contentType string) {
	name := strings.ToLower(filename)
	if declaredType != "" { mediaType = declaredType }
	switch {
	case strings.Contains(name, ".jpg") || strings.Contains(name, ".jpeg") || strings.Contains(name, ".png") || strings.Contains(name, ".webp") || strings.Contains(name, ".gif"):
		if mediaType == "" { mediaType = "image" }
		contentType = headerCT
		if contentType == "" { contentType = "application/octet-stream" }
	case strings.Contains(name, ".mp3") || strings.Contains(name, ".ogg") || strings.Contains(name, ".m4a") || strings.Contains(name, ".wav") || strings.Contains(name, ".aac"):
		if mediaType == "" { mediaType = "audio" }
		contentType = headerCT
	case strings.Contains(name, ".mp4") || strings.Contains(name, ".3gp") || strings.Contains(name, ".mov") || strings.Contains(name, ".mkv") || strings.Contains(name, ".avi"):
		if mediaType == "" { mediaType = "video" }
		contentType = headerCT
	default:
		if mediaType == "" { mediaType = "document" }
		contentType = headerCT
	}
	if contentType == "" { contentType = "application/octet-stream" }
	return
}

func (u *MediaUploader) UploadAndSend(ctx context.Context, conv *models.Conversation, file *multipart.FileHeader, declaredType, caption string) (*models.Message, error) {
	f, err := file.Open()
	if err != nil { return nil, err }
	defer f.Close()

	mediaType, contentType := detectType(file.Filename, file.Header.Get("Content-Type"), declaredType)
	objectPath := storage.Join("whatsapp-crm", mediaType, time.Now().Format("2006/01/02"), uuid.New().String()+filepath.Ext(file.Filename))

	// Save to storage (private)
	savedPath, err := u.store.Save(ctx, f.(io.Reader), objectPath, contentType)
	if err != nil { return nil, fmt.Errorf("store: %w", err) }

	// Signed URL
	signedURL, err := u.store.SignedURL(ctx, savedPath, u.expiry)
	if err != nil { return nil, fmt.Errorf("signedurl: %w", err) }

	// Send via WA
	var resp *whatsapp.SendMessageResponse
	switch mediaType {
	case "image":
		resp, err = u.wa.SendImageMessage(conv.Customer.WhatsAppID, signedURL, caption)
	case "document":
		resp, err = u.wa.SendDocumentMessage(conv.Customer.WhatsAppID, signedURL, file.Filename, caption)
	case "audio":
		resp, err = u.wa.SendAudioMessage(conv.Customer.WhatsAppID, signedURL)
	case "video":
		resp, err = u.wa.SendVideoMessage(conv.Customer.WhatsAppID, signedURL, caption)
	default:
		resp, err = u.wa.SendDocumentMessage(conv.Customer.WhatsAppID, signedURL, file.Filename, caption)
	}
	if err != nil { return nil, err }

	now := time.Now()
	msg := models.Message{ConversationID: conv.ID, WhatsAppID: resp.ID, Type: models.MessageType(mediaType), Direction: models.MessageDirectionOutbound, Status: models.MessageStatusSent, MediaURL: signedURL, Caption: caption, FileName: file.Filename, SentAt: &now}
	return &msg, nil
}
