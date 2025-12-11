package onesignal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yousoon/apps/services/notification-service/internal/domain"
)

// Client pour l'API OneSignal
type Client struct {
	appID      string
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewClient crée un nouveau client OneSignal
func NewClient(appID, apiKey string) *Client {
	return &Client{
		appID:  appID,
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://onesignal.com/api/v1",
	}
}

// NotificationRequest représente une requête de notification
type NotificationRequest struct {
	AppID            string                 `json:"app_id"`
	IncludePlayerIDs []string               `json:"include_player_ids,omitempty"`
	IncludeAliases   *IncludeAliases        `json:"include_aliases,omitempty"`
	TargetChannel    string                 `json:"target_channel,omitempty"`
	Headings         map[string]string      `json:"headings"`
	Contents         map[string]string      `json:"contents"`
	BigPicture       string                 `json:"big_picture,omitempty"`
	IOSAttachments   map[string]string      `json:"ios_attachments,omitempty"`
	Data             map[string]interface{} `json:"data,omitempty"`
	URL              string                 `json:"url,omitempty"`
	IOSBadgeType     string                 `json:"ios_badgeType,omitempty"`
	IOSBadgeCount    int                    `json:"ios_badgeCount,omitempty"`
	TTL              int                    `json:"ttl,omitempty"`
	Priority         int                    `json:"priority,omitempty"`
}

// IncludeAliases pour cibler par external_id
type IncludeAliases struct {
	ExternalID []string `json:"external_id"`
}

// NotificationResponse réponse de l'API
type NotificationResponse struct {
	ID         string   `json:"id"`
	Recipients int      `json:"recipients"`
	Errors     []string `json:"errors,omitempty"`
}

// SendPush envoie une notification push
func (c *Client) SendPush(ctx context.Context, notification *domain.Notification, tokens []string) error {
	// Préparer la data
	data := map[string]interface{}{
		"type":        notification.Type().String(),
		"relatedType": notification.RelatedType(),
		"relatedId":   notification.RelatedID(),
	}

	req := NotificationRequest{
		AppID:            c.appID,
		IncludePlayerIDs: tokens,
		Headings: map[string]string{
			"en": notification.Title(),
			"fr": notification.Title(),
		},
		Contents: map[string]string{
			"en": notification.Body(),
			"fr": notification.Body(),
		},
		Data:     data,
		TTL:      86400, // 24 heures
		Priority: 10,
	}

	if imageURL := notification.ImageURL(); imageURL != nil && *imageURL != "" {
		req.BigPicture = *imageURL
		req.IOSAttachments = map[string]string{
			"id": *imageURL,
		}
	}

	return c.sendRequest(ctx, "/notifications", req)
}

// SendToUser envoie une notification à un utilisateur via son external_id
func (c *Client) SendToUser(ctx context.Context, notification *domain.Notification, userID string) error {
	data := map[string]interface{}{
		"type":        notification.Type().String(),
		"relatedType": notification.RelatedType(),
		"relatedId":   notification.RelatedID(),
	}

	req := NotificationRequest{
		AppID: c.appID,
		IncludeAliases: &IncludeAliases{
			ExternalID: []string{userID},
		},
		TargetChannel: "push",
		Headings: map[string]string{
			"en": notification.Title(),
			"fr": notification.Title(),
		},
		Contents: map[string]string{
			"en": notification.Body(),
			"fr": notification.Body(),
		},
		Data:     data,
		TTL:      86400,
		Priority: 10,
	}

	if imageURL := notification.ImageURL(); imageURL != nil && *imageURL != "" {
		req.BigPicture = *imageURL
		req.IOSAttachments = map[string]string{
			"id": *imageURL,
		}
	}

	return c.sendRequest(ctx, "/notifications", req)
}

func (c *Client) sendRequest(ctx context.Context, path string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp NotificationResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("onesignal error: %v", errResp.Errors)
	}

	return nil
}
