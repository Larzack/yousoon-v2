package domain

import "time"

// Entity is the base interface for all domain entities.
// Entities have identity that persists over time.
type Entity interface {
	// GetID returns the entity's unique identifier as a string
	GetID() string
}

// BaseEntity provides common fields for all entities.
type BaseEntity struct {
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// NewBaseEntity creates a new base entity with timestamps.
func NewBaseEntity() BaseEntity {
	now := time.Now().UTC()
	return BaseEntity{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// MarkUpdated updates the UpdatedAt timestamp.
func (e *BaseEntity) MarkUpdated() {
	e.UpdatedAt = time.Now().UTC()
}

// MarkDeleted performs a soft delete by setting DeletedAt.
func (e *BaseEntity) MarkDeleted() {
	now := time.Now().UTC()
	e.DeletedAt = &now
	e.UpdatedAt = now
}

// IsDeleted returns true if the entity has been soft deleted.
func (e *BaseEntity) IsDeleted() bool {
	return e.DeletedAt != nil
}

// Restore undoes a soft delete.
func (e *BaseEntity) Restore() {
	e.DeletedAt = nil
	e.UpdatedAt = time.Now().UTC()
}

// Timestamps returns the entity's timestamps.
func (e *BaseEntity) Timestamps() (created, updated time.Time, deleted *time.Time) {
	return e.CreatedAt, e.UpdatedAt, e.DeletedAt
}

// =============================================================================
// Auditable Entity
// =============================================================================

// AuditInfo contains audit information for an entity.
type AuditInfo struct {
	CreatedBy string     `json:"created_by,omitempty" bson:"created_by,omitempty"`
	UpdatedBy string     `json:"updated_by,omitempty" bson:"updated_by,omitempty"`
	DeletedBy string     `json:"deleted_by,omitempty" bson:"deleted_by,omitempty"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// NewAuditInfo creates audit info with the creator.
func NewAuditInfo(createdBy string) AuditInfo {
	now := time.Now().UTC()
	return AuditInfo{
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// MarkUpdatedBy updates the entity with the modifier.
func (a *AuditInfo) MarkUpdatedBy(updatedBy string) {
	a.UpdatedBy = updatedBy
	a.UpdatedAt = time.Now().UTC()
}

// MarkDeletedBy performs a soft delete with the deleter.
func (a *AuditInfo) MarkDeletedBy(deletedBy string) {
	now := time.Now().UTC()
	a.DeletedBy = deletedBy
	a.DeletedAt = &now
	a.UpdatedAt = now
}

// =============================================================================
// Versioned Entity (Optimistic Locking)
// =============================================================================

// VersionedEntity adds version tracking for optimistic locking.
type VersionedEntity struct {
	BaseEntity
	Version int64 `json:"version" bson:"version"`
}

// NewVersionedEntity creates a new versioned entity.
func NewVersionedEntity() VersionedEntity {
	return VersionedEntity{
		BaseEntity: NewBaseEntity(),
		Version:    1,
	}
}

// IncrementVersion increments the version for optimistic locking.
func (e *VersionedEntity) IncrementVersion() {
	e.Version++
	e.MarkUpdated()
}

// GetVersion returns the current version.
func (e *VersionedEntity) GetVersion() int64 {
	return e.Version
}
