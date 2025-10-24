package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
	"whatsapp-crm/internal/config"
)

type Client struct {
	APIURL   string
	APIToken string
	client   *http.Client
}

type SendMessageRequest struct {
	To      string `json:"to"`
	Type    string `json:"type"`
	Message interface{} `json:"message"`
}

type TextMessage struct {
	Body string `json:"body"`
}

type MediaMessage struct {
	URL      string `json:"url,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename,omitempty"`
}

type TemplateMessage struct {
	Name       string                 `json:"name"`
	Language   TemplateLanguage       `json:"language"`
	Components []TemplateComponent    `json:"components,omitempty"`
}

type TemplateLanguage struct {
	Code string `json:"code"`
}

type TemplateComponent struct {
	Type       string                    `json:"type"`
	Parameters []TemplateParameter      `json:"parameters,omitempty"`
}

type TemplateParameter struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type SendMessageResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
	Message  string `json:"message,omitempty"`
}

type MessageStatus struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	To        string    `json:"to"`
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		APIURL:   cfg.WhatsAppAPIURL,
		APIToken: cfg.WhatsAppAPIToken,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendTextMessage sends a text message
func (c *Client) SendTextMessage(to, message string) (*SendMessageResponse, error) {
	req := SendMessageRequest{
		To:   to,
		Type: "text",
		Message: TextMessage{
			Body: message,
		},
	}

	return c.sendMessage(req)
}

// SendImageMessage sends an image message
func (c *Client) SendImageMessage(to, imageURL, caption string) (*SendMessageResponse, error) {
	req := SendMessageRequest{
		To:   to,
		Type: "image",
		Message: MediaMessage{
			URL:     imageURL,
			Caption: caption,
		},
	}

	return c.sendMessage(req)
}

// SendDocumentMessage sends a document message
func (c *Client) SendDocumentMessage(to, documentURL, filename, caption string) (*SendMessageResponse, error) {
	req := SendMessageRequest{
		To:   to,
		Type: "document",
		Message: MediaMessage{
			URL:      documentURL,
			Filename: filename,
			Caption:  caption,
		},
	}

	return c.sendMessage(req)
}

// SendTemplateMessage sends a template message
func (c *Client) SendTemplateMessage(to, templateName, languageCode string, components []TemplateComponent) (*SendMessageResponse, error) {
	req := SendMessageRequest{
		To:   to,
		Type: "template",
		Message: TemplateMessage{
			Name: templateName,
			Language: TemplateLanguage{
				Code: languageCode,
			},
			Components: components,
		},
	}

	return c.sendMessage(req)
}

func (c *Client) sendMessage(req SendMessageRequest) (*SendMessageResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.APIURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIToken)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response SendMessageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &response, fmt.Errorf("API error: %s", response.Error)
	}

	return &response, nil
}

// UploadMedia uploads media file and returns URL
func (c *Client) UploadMedia(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", c.APIURL+"/media", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.APIToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var uploadResp struct {
		URL   string `json:"url"`
		Error string `json:"error"`
	}

	if err := json.Unmarshal(respBody, &uploadResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", uploadResp.Error)
	}

	return uploadResp.URL, nil
}

// GetMessageStatus gets message delivery status
func (c *Client) GetMessageStatus(messageID string) (*MessageStatus, error) {
	req, err := http.NewRequest("GET", c.APIURL+"/messages/"+messageID+"/status", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var status MessageStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &status, nil
}