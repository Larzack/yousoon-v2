package domain

import (
	"time"
)

// =============================================================================
// Domain Events
// =============================================================================

// PartnerEvent is the base interface for all partner domain events.
type PartnerEvent interface {
	EventName() string
	OccurredAt() time.Time
	AggregateID() string
}

// =============================================================================
// Partner Registration Events
// =============================================================================

// PartnerRegisteredEvent is published when a new partner registers.
type PartnerRegisteredEvent struct {
	PartnerID   PartnerID `json:"partnerId"`
	OwnerUserID UserID    `json:"ownerUserId"`
	CompanyName string    `json:"companyName"`
	Timestamp   time.Time `json:"timestamp"`
}

// NewPartnerRegisteredEvent creates a new PartnerRegisteredEvent.
func NewPartnerRegisteredEvent(partnerID PartnerID, ownerUserID UserID, companyName string) PartnerRegisteredEvent {
	return PartnerRegisteredEvent{
		PartnerID:   partnerID,
		OwnerUserID: ownerUserID,
		CompanyName: companyName,
		Timestamp:   time.Now(),
	}
}

// EventName returns the event name.
func (e PartnerRegisteredEvent) EventName() string {
	return "partner.registered"
}

// OccurredAt returns when the event occurred.
func (e PartnerRegisteredEvent) OccurredAt() time.Time {
	return e.Timestamp
}

// AggregateID returns the aggregate ID.
func (e PartnerRegisteredEvent) AggregateID() string {
	return e.PartnerID.String()
}

// =============================================================================
// Partner Verification Events
// =============================================================================

// PartnerVerifiedEvent is published when a partner is verified.
type PartnerVerifiedEvent struct {
	PartnerID  PartnerID `json:"partnerId"`
	VerifiedBy string    `json:"verifiedBy"`
	Timestamp  time.Time `json:"timestamp"`
}

// NewPartnerVerifiedEvent creates a new PartnerVerifiedEvent.
func NewPartnerVerifiedEvent(partnerID PartnerID, verifiedBy string) PartnerVerifiedEvent {
	return PartnerVerifiedEvent{
		PartnerID:  partnerID,
		VerifiedBy: verifiedBy,
		Timestamp:  time.Now(),
	}
}

// EventName returns the event name.
func (e PartnerVerifiedEvent) EventName() string {
	return "partner.verified"
}

// OccurredAt returns when the event occurred.
func (e PartnerVerifiedEvent) OccurredAt() time.Time {
	return e.Timestamp
}

// AggregateID returns the aggregate ID.
func (e PartnerVerifiedEvent) AggregateID() string {
	return e.PartnerID.String()
}

// =============================================================================
// Partner Status Events
// =============================================================================

// PartnerSuspendedEvent is published when a partner is suspended.
type PartnerSuspendedEvent struct {
	PartnerID PartnerID `json:"partnerId"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

// NewPartnerSuspendedEvent creates a new PartnerSuspendedEvent.
func NewPartnerSuspendedEvent(partnerID PartnerID, reason string) PartnerSuspendedEvent {
	return PartnerSuspendedEvent{
		PartnerID: partnerID,
		Reason:    reason,
		Timestamp: time.Now(),
	}
}

// EventName returns the event name.
func (e PartnerSuspendedEvent) EventName() string {
	return "partner.suspended"
}

// OccurredAt returns when the event occurred.
func (e PartnerSuspendedEvent) OccurredAt() time.Time {
	return e.Timestamp
}

// AggregateID returns the aggregate ID.
func (e PartnerSuspendedEvent) AggregateID() string {
	return e.PartnerID.String()
}

// =============================================================================
// Establishment Events
// =============================================================================

// EstablishmentAddedEvent is published when an establishment is added.
type EstablishmentAddedEvent struct {
	PartnerID       PartnerID       `json:"partnerId"`
	EstablishmentID EstablishmentID `json:"establishmentId"`
	Name            string          `json:"name"`
	Location        GeoLocation     `json:"location"`
	Timestamp       time.Time       `json:"timestamp"`
}

// NewEstablishmentAddedEvent creates a new EstablishmentAddedEvent.
func NewEstablishmentAddedEvent(partnerID PartnerID, estID EstablishmentID, name string, location GeoLocation) EstablishmentAddedEvent {
	return EstablishmentAddedEvent{
		PartnerID:       partnerID,
		EstablishmentID: estID,
		Name:            name,
		Location:        location,
		Timestamp:       time.Now(),
	}
}

// EventName returns the event name.
func (e EstablishmentAddedEvent) EventName() string {
	return "partner.establishment_added"
}

// OccurredAt returns when the event occurred.
func (e EstablishmentAddedEvent) OccurredAt() time.Time {
	return e.Timestamp
}

// AggregateID returns the aggregate ID.
func (e EstablishmentAddedEvent) AggregateID() string {
	return e.PartnerID.String()
}

// =============================================================================
// Team Member Events
// =============================================================================

// TeamMemberInvitedEvent is published when a team member is invited.
type TeamMemberInvitedEvent struct {
	PartnerID PartnerID `json:"partnerId"`
	Email     Email     `json:"email"`
	Role      TeamRole  `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}

// NewTeamMemberInvitedEvent creates a new TeamMemberInvitedEvent.
func NewTeamMemberInvitedEvent(partnerID PartnerID, email Email, role TeamRole) TeamMemberInvitedEvent {
	return TeamMemberInvitedEvent{
		PartnerID: partnerID,
		Email:     email,
		Role:      role,
		Timestamp: time.Now(),
	}
}

// EventName returns the event name.
func (e TeamMemberInvitedEvent) EventName() string {
	return "partner.team_member_invited"
}

// OccurredAt returns when the event occurred.
func (e TeamMemberInvitedEvent) OccurredAt() time.Time {
	return e.Timestamp
}

// AggregateID returns the aggregate ID.
func (e TeamMemberInvitedEvent) AggregateID() string {
	return e.PartnerID.String()
}

// TeamMemberJoinedEvent is published when a team member accepts an invitation.
type TeamMemberJoinedEvent struct {
	PartnerID    PartnerID    `json:"partnerId"`
	TeamMemberID TeamMemberID `json:"teamMemberId"`
	UserID       UserID       `json:"userId"`
	Timestamp    time.Time    `json:"timestamp"`
}

// NewTeamMemberJoinedEvent creates a new TeamMemberJoinedEvent.
func NewTeamMemberJoinedEvent(partnerID PartnerID, memberID TeamMemberID, userID UserID) TeamMemberJoinedEvent {
	return TeamMemberJoinedEvent{
		PartnerID:    partnerID,
		TeamMemberID: memberID,
		UserID:       userID,
		Timestamp:    time.Now(),
	}
}

// EventName returns the event name.
func (e TeamMemberJoinedEvent) EventName() string {
	return "partner.team_member_joined"
}

// OccurredAt returns when the event occurred.
func (e TeamMemberJoinedEvent) OccurredAt() time.Time {
	return e.Timestamp
}

// AggregateID returns the aggregate ID.
func (e TeamMemberJoinedEvent) AggregateID() string {
	return e.PartnerID.String()
}

// =============================================================================
// Statistics Events
// =============================================================================

// PartnerStatsUpdatedEvent is published when partner stats are updated.
type PartnerStatsUpdatedEvent struct {
	PartnerID PartnerID    `json:"partnerId"`
	Stats     PartnerStats `json:"stats"`
	Timestamp time.Time    `json:"timestamp"`
}

// NewPartnerStatsUpdatedEvent creates a new PartnerStatsUpdatedEvent.
func NewPartnerStatsUpdatedEvent(partnerID PartnerID, stats PartnerStats) PartnerStatsUpdatedEvent {
	return PartnerStatsUpdatedEvent{
		PartnerID: partnerID,
		Stats:     stats,
		Timestamp: time.Now(),
	}
}

// EventName returns the event name.
func (e PartnerStatsUpdatedEvent) EventName() string {
	return "partner.stats_updated"
}

// OccurredAt returns when the event occurred.
func (e PartnerStatsUpdatedEvent) OccurredAt() time.Time {
	return e.Timestamp
}

// AggregateID returns the aggregate ID.
func (e PartnerStatsUpdatedEvent) AggregateID() string {
	return e.PartnerID.String()
}
