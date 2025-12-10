package nats

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/yousoon/apps/services/notification-service/internal/domain"
)

// EventHandler gère les événements NATS pour les notifications
type EventHandler struct {
	js                  nats.JetStreamContext
	notificationRepo    domain.NotificationRepository
	deviceTokenRepo     domain.DeviceTokenRepository
	templateRepo        domain.TemplateRepository
	notificationService *NotificationService
}

// NotificationService service d'envoi de notifications
type NotificationService struct {
	pushSender  PushSender
	emailSender EmailSender
	smsSender   SMSSender
}

// PushSender interface pour l'envoi de push
type PushSender interface {
	SendToUser(ctx context.Context, notification *domain.Notification, userID string) error
}

// EmailSender interface pour l'envoi d'email
type EmailSender interface {
	SendEmail(ctx context.Context, notification *domain.Notification, email string) error
}

// SMSSender interface pour l'envoi de SMS
type SMSSender interface {
	SendSMS(ctx context.Context, notification *domain.Notification, phone string) error
}

// NewEventHandler crée un nouveau handler d'événements
func NewEventHandler(
	js nats.JetStreamContext,
	notificationRepo domain.NotificationRepository,
	deviceTokenRepo domain.DeviceTokenRepository,
	templateRepo domain.TemplateRepository,
	pushSender PushSender,
	emailSender EmailSender,
	smsSender SMSSender,
) *EventHandler {
	return &EventHandler{
		js:               js,
		notificationRepo: notificationRepo,
		deviceTokenRepo:  deviceTokenRepo,
		templateRepo:     templateRepo,
		notificationService: &NotificationService{
			pushSender:  pushSender,
			emailSender: emailSender,
			smsSender:   smsSender,
		},
	}
}

// Subscribe s'abonne aux événements
func (h *EventHandler) Subscribe(ctx context.Context) error {
	// S'abonner aux événements de booking
	if _, err := h.js.Subscribe("yousoon.events.booking.*", h.handleBookingEvents, nats.Durable("notification-booking")); err != nil {
		return err
	}

	// S'abonner aux événements d'offres
	if _, err := h.js.Subscribe("yousoon.events.offer.*", h.handleOfferEvents, nats.Durable("notification-offer")); err != nil {
		return err
	}

	// S'abonner aux événements utilisateur
	if _, err := h.js.Subscribe("yousoon.events.user.*", h.handleUserEvents, nats.Durable("notification-user")); err != nil {
		return err
	}

	log.Println("Notification event handlers subscribed")
	return nil
}

func (h *EventHandler) handleBookingEvents(msg *nats.Msg) {
	ctx := context.Background()

	var event struct {
		Type      string `json:"type"`
		UserID    string `json:"userId"`
		OutingID  string `json:"outingId"`
		OfferID   string `json:"offerId"`
		PartnerID string `json:"partnerId"`
	}

	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("Error unmarshaling booking event: %v", err)
		msg.Nak()
		return
	}

	var notification *domain.Notification

	switch event.Type {
	case "OutingCreated":
		notification = &domain.Notification{
			UserID:      event.UserID,
			Type:        domain.NotificationTypeBookingConfirmed,
			Channel:     domain.NotificationChannelPush,
			Title:       "Réservation confirmée !",
			Body:        "Votre sortie a été réservée. Présentez le QR code pour valider.",
			RelatedType: "outing",
			RelatedID:   event.OutingID,
			Status:      domain.NotificationStatusPending,
		}

	case "OutingCheckedIn":
		notification = &domain.Notification{
			UserID:      event.UserID,
			Type:        domain.NotificationTypeCheckInSuccess,
			Channel:     domain.NotificationChannelPush,
			Title:       "Check-in réussi !",
			Body:        "Profitez bien de votre sortie !",
			RelatedType: "outing",
			RelatedID:   event.OutingID,
			Status:      domain.NotificationStatusPending,
		}

	case "OutingCancelled":
		notification = &domain.Notification{
			UserID:      event.UserID,
			Type:        domain.NotificationTypeBookingCancelled,
			Channel:     domain.NotificationChannelPush,
			Title:       "Réservation annulée",
			Body:        "Votre réservation a été annulée.",
			RelatedType: "outing",
			RelatedID:   event.OutingID,
			Status:      domain.NotificationStatusPending,
		}
	}

	if notification != nil {
		if err := h.sendNotification(ctx, notification); err != nil {
			log.Printf("Error sending notification: %v", err)
			msg.Nak()
			return
		}
	}

	msg.Ack()
}

func (h *EventHandler) handleOfferEvents(msg *nats.Msg) {
	ctx := context.Background()

	var event struct {
		Type    string  `json:"type"`
		OfferID string  `json:"offerId"`
		Lat     float64 `json:"lat"`
		Lng     float64 `json:"lng"`
	}

	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("Error unmarshaling offer event: %v", err)
		msg.Nak()
		return
	}

	switch event.Type {
	case "OfferPublished":
		// TODO: Envoyer des notifications aux utilisateurs à proximité
		// Cela nécessiterait de récupérer les utilisateurs par géolocalisation
		log.Printf("Offer published: %s at (%f, %f)", event.OfferID, event.Lat, event.Lng)
	}

	msg.Ack()
}

func (h *EventHandler) handleUserEvents(msg *nats.Msg) {
	ctx := context.Background()

	var event struct {
		Type   string `json:"type"`
		UserID string `json:"userId"`
		Email  string `json:"email"`
	}

	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("Error unmarshaling user event: %v", err)
		msg.Nak()
		return
	}

	switch event.Type {
	case "UserRegistered":
		// Envoyer un email de bienvenue
		notification := &domain.Notification{
			UserID:  event.UserID,
			Type:    domain.NotificationTypeSystem,
			Channel: domain.NotificationChannelEmail,
			Title:   "Bienvenue sur Yousoon !",
			Body:    "Découvrez les meilleures sorties à prix réduit près de chez vous.",
			Status:  domain.NotificationStatusPending,
		}

		if err := h.sendNotification(ctx, notification); err != nil {
			log.Printf("Error sending welcome email: %v", err)
		}

	case "UserIdentityVerified":
		notification := &domain.Notification{
			UserID:  event.UserID,
			Type:    domain.NotificationTypeSystem,
			Channel: domain.NotificationChannelPush,
			Title:   "Identité vérifiée !",
			Body:    "Votre identité a été vérifiée. Vous pouvez maintenant réserver des sorties.",
			Status:  domain.NotificationStatusPending,
		}

		if err := h.sendNotification(ctx, notification); err != nil {
			log.Printf("Error sending verification notification: %v", err)
		}
	}

	msg.Ack()
}

func (h *EventHandler) sendNotification(ctx context.Context, notification *domain.Notification) error {
	// Sauvegarder la notification
	if err := h.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	// Envoyer selon le channel
	var err error
	switch notification.Channel {
	case domain.NotificationChannelPush:
		err = h.notificationService.pushSender.SendToUser(ctx, notification, notification.UserID)
	case domain.NotificationChannelEmail:
		// TODO: Récupérer l'email de l'utilisateur via gRPC
		// err = h.notificationService.emailSender.SendEmail(ctx, notification, userEmail)
	case domain.NotificationChannelSMS:
		// TODO: Récupérer le téléphone de l'utilisateur via gRPC
		// err = h.notificationService.smsSender.SendSMS(ctx, notification, userPhone)
	}

	if err != nil {
		notification.Status = domain.NotificationStatusFailed
		notification.Error = err.Error()
	} else {
		notification.Status = domain.NotificationStatusSent
	}

	return h.notificationRepo.Update(ctx, notification)
}
