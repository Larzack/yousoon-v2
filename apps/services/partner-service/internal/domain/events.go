package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// Domain Events
// =============================================================================

// =============================================================================
// Partner Registration Events
// =============================================================================

// PartnerRegisteredEvent is published when a new partner registers.
type PartnerRegisteredEvent struct {
	ID          string    `json:"event_id"`
	PartnerID   PartnerID `json:"partnerId"`
	OwnerUserID UserID    `json:"ownerUserId"`
	CompanyName string    `json:"companyName"`
	Timestamp   time.Time `json:"timestamp"`
}

// NewPartnerRegisteredEvent creates a new PartnerRegisteredEvent.
func NewPartnerRegisteredEvent(partnerID PartnerID, ownerUserID UserID, companyName string) PartnerRegisteredEvent {
	return PartnerRegisteredEvent{
		ID:          uuid.New().String(),
		PartnerID:   partnerID,
		OwnerUserID: ownerUserID,
		CompanyName: companyName,
		Timestamp:   time.Now(),
	}
}

func (e PartnerRegisteredEvent) EventID() string          { return e.ID }
func (e PartnerRegisteredEvent) EventName() string        { return "partner.registered" }
func (e PartnerRegisteredEvent) OccurredAt() time.Time    { return e.Timestamp }
func (e PartnerRegisteredEvent) AggregateID() string      { return e.PartnerID.String() }
func (e PartnerRegisteredEvent) AggregateType() string    { return "Partner" }
func (e PartnerRegisteredEvent) Version() int             { return 1 }
func (e PartnerRegisteredEvent) Payload() ([]byte, error) { return json.Marshal(e) }

// =============================================================================
// Partner Verification Events
// =============================================================================

// PartnerVerifiedEvent is published when a partner is verified.
type PartnerVerifiedEvent struct {
	ID         string    `json:"event_id"`
	PartnerID  PartnerID `json:"partnerId"`
	VerifiedBy string    `json:"verifiedBy"`
	Timestamp  time.Time `json:"timestamp"`
}

// NewPartnerVerifiedEvent creates a new PartnerVerifiedEvent.
func NewPartnerVerifiedEvent(partnerID PartnerID, verifiedBy string) PartnerVerifiedEvent {
	return PartnerVerifiedEvent{
		ID:         uuid.New().String(),
		PartnerID:  partnerID,
		VerifiedBy: verifiedBy,
		Timestamp:  time.Now(),
	}
}

func (e PartnerVerifiedEvent) EventID() string          { return e.ID }
func (e PartnerVerifiedEvent) EventName() string        { return "partner.verified" }
func (e PartnerVerifiedEvent) OccurredAt() time.Time    { return e.Timestamp }
func (e PartnerVerifiedEvent) AggregateID() string      { return e.PartnerID.String() }
func (e PartnerVerifiedEvent) AggregateType() string    { return "Partner" }
func (e PartnerVerifiedEvent) Version() int             { return 1 }
func (e PartnerVerifiedEvent) Payload() ([]byte, error) { return json.Marshal(e) }

// =============================================================================
// Partner Status Events
// =============================================================================

// PartnerSuspendedEvent is published when a partner is suspended.
type PartnerSuspendedEvent struct {
	ID        string    `json:"event_id"`
	PartnerID PartnerID `json:"partnerId"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

// NewPartnerSuspendedEvent creates a new PartnerSuspendedEvent.
func NewPartnerSuspendedEvent(partnerID PartnerID, reason string) PartnerSuspendedEvent {
	return PartnerSuspendedEvent{
		ID:        uuid.New().String(),
		PartnerID: partnerID,
		Reason:    reason,
		Timestamp: time.Now(),
	}
}

func (e PartnerSuspendedEvent) EventID() string          { return e.ID }
func (e PartnerSuspendedEvent) EventName() string        { return "partner.suspended" }
func (e PartnerSuspendedEvent) OccurredAt() time.Time    { return e.Timestamp }
func (e PartnerSuspendedEvent) AggregateID() string      { return e.PartnerID.String() }
func (e PartnerSuspendedEvent) AggregateType() string    { return "Partner" }
func (e PartnerSuspendedEvent) Version() int             { return 1 }
func (e PartnerSuspendedEvent) Payload() ([]byte, error) { return json.Marshal(e) }

// =============================================================================
// Establishment Events
// =============================================================================

// EstablishmentAddedEvent is published when an establishment is added.
type EstablishmentAddedEvent struct {
	ID              string          `json:"event_id"`
	PartnerID       PartnerID       `json:"partnerId"`
	EstablishmentID EstablishmentID `json:"establishmentId"`
	Name            string          `json:"name"`
	Location        GeoLocation     `json:"location"`
	Timestamp       time.Time       `json:"timestamp"`
}

// NewEstablishmentAddedEvent creates a new EstablishmentAddedEvent.
func NewEstablishmentAddedEvent(partnerID PartnerID, estID EstablishmentID, name string, location GeoLocation) EstablishmentAddedEvent {
	return EstablishmentAddedEvent{
		ID:              uuid.New().String(),
		PartnerID:       partnerID,
		EstablishmentID: estID,
		Name:            name,
		Location:        location,
		Timestamp:       time.Now(),
	}
}

func (e EstablishmentAddedEvent) EventID() string          { return e.ID }
func (e EstablishmentAddedEvent) EventName() string        { return "partner.establishment_added" }
func (e EstablishmentAddedEvent) OccurredAt() time.Time    { return e.Timestamp }
func (e EstablishmentAddedEvent) AggregateID() string      { return e.PartnerID.String() }
func (e EstablishmentAddedEvent) AggregateType() string    { return "Partner" }
func (e EstablishmentAddedEvent) Version() int             { return 1 }
func (e EstablishmentAddedEvent) Payload() ([]byte, error) { return json.Marshal(e) }

// =============================================================================
// Team Member Events
// =============================================================================

// TeamMemberInvitedEvent is published when a team member is invited.
type TeamMemberInvitedEvent struct {
	ID        string    `json:"event_id"`
	PartnerID PartnerID `json:"partnerId"`
	Email     Email     `json:"email"`
	Role      TeamRole  `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}

// NewTeamMemberInvitedEvent creates a new TeamMemberInvitedEvent.
func NewTeamMemberInvitedEvent(partnerID PartnerID, email Email, role TeamRole) TeamMemberInvitedEvent {
	return TeamMemberInvitedEvent{
		ID:        uuid.New().String(),
		PartnerID: partnerID,
		Email:     email,
		Role:      role,
		Timestamp: time.Now(),
	}
}

func (e TeamMemberInvitedEvent) EventID() string          { return e.ID }
func (e TeamMemberInvitedEvent) EventName() string        { return "partner.team_member_invited" }
func (e TeamMemberInvitedEvent) OccurredAt() time.Time    { return e.Timestamp }
func (e TeamMemberInvitedEvent) AggregateID() string      { return e.PartnerID.String() }
func (e TeamMemberInvitedEvent) AggregateType() string    { return "Partner" }
func (e TeamMemberInvitedEvent) Version() int             { return 1 }
func (e TeamMemberInvitedEvent) Payload() ([]byte, error) { return json.Marshal(e) }

// TeamMemberJoinedEvent is published when a team member accepts an invitation.
type TeamMemberJoinedEvent struct {
	ID           string       `json:"event_id"`
	PartnerID    PartnerID    `json:"partnerId"`
	TeamMemberID TeamMemberID `json:"teamMemberId"`
	UserID       UserID       `json:"userId"`
	Timestamp    time.Time    `json:"timestamp"`
}

// NewTeamMemberJoinedEvent creates a new TeamMemberJoinedEvent.
func NewTeamMemberJoinedEvent(partnerID PartnerID, memberID TeamMemberID, userID UserID) TeamMemberJoinedEvent {
	return TeamMemberJoinedEvent{
		ID:           uuid.New().String(),
		PartnerID:    partnerID,
		TeamMemberID: memberID,
		UserID:       userID,
		Timestamp:    time.Now(),
	}
}

func (e TeamMemberJoinedEvent) EventID() string          { return e.ID }
func (e TeamMemberJoinedEvent) EventName() string        { return "partner.team_member_joined" }
func (e TeamMemberJoinedEvent) OccurredAt() time.Time    { return e.Timestamp }
func (e TeamMemberJoinedEvent) AggregateID() string      { return e.PartnerID.String() }
func (e TeamMemberJoinedEvent) AggregateType() string    { return "Partner" }
func (e TeamMemberJoinedEvent) Version() int             { return 1 }
func (e TeamMemberJoinedEvent) Payload() ([]byte, error) { return json.Marshal(e) }

// =============================================================================
// Statistics Events
// =============================================================================

// PartnerStatsUpdatedEvent is published when partner stats are updated.
type PartnerStatsUpdatedEvent struct {
	ID        string       `json:"event_id"`
	PartnerID PartnerID    `json:"partnerId"`
	Stats     PartnerStats `json:"stats"`
	Timestamp time.Time    `json:"timestamp"`
}

// NewPartnerStatsUpdatedEvent creates a new PartnerStatsUpdatedEvent.
func NewPartnerStatsUpdatedEvent(partnerID PartnerID, stats PartnerStats) PartnerStatsUpdatedEvent {
	return PartnerStatsUpdatedEvent{
		ID:        uuid.New().String(),
		PartnerID: partnerID,
		Stats:     stats,
		Timestamp: time.Now(),
	}
}

func (e PartnerStatsUpdatedEvent) EventID() string          { return e.ID }
func (e PartnerStatsUpdatedEvent) EventName() string        { return "partner.stats_updated" }
func (e PartnerStatsUpdatedEvent) OccurredAt() time.Time    { return e.Timestamp }
func (e PartnerStatsUpdatedEvent) AggregateID() string      { return e.PartnerID.String() }
func (e PartnerStatsUpdatedEvent) AggregateType() string    { return "Partner" }
func (e PartnerStatsUpdatedEvent) Version() int             { return 1 }
func (e PartnerStatsUpdatedEvent) Payload() ([]byte, error) { return json.Marshal(e) }
