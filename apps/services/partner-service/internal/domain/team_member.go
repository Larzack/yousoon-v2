package domain

import (
	"time"
)

// =============================================================================
// TeamMember Entity
// =============================================================================

// TeamMember represents a person in a partner's team.
// It's an entity within the Partner aggregate.
type TeamMember struct {
	// Identity
	ID        TeamMemberID `json:"id" bson:"id"`
	PartnerID PartnerID    `json:"partnerId" bson:"partnerId"`

	// User Link (optional until invitation is accepted)
	UserID *UserID `json:"userId,omitempty" bson:"userId,omitempty"`

	// Contact Info
	Email     Email  `json:"email" bson:"email"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`

	// Role
	Role TeamRole `json:"role" bson:"role"`

	// Status
	Status TeamMemberStatus `json:"status" bson:"status"`

	// Timestamps
	InvitedAt time.Time  `json:"invitedAt" bson:"invitedAt"`
	JoinedAt  *time.Time `json:"joinedAt,omitempty" bson:"joinedAt,omitempty"`
}

// =============================================================================
// TeamMember Factory
// =============================================================================

// NewTeamMember creates a new team member invitation.
func NewTeamMember(email Email, firstName, lastName string, role TeamRole) TeamMember {
	return TeamMember{
		ID:        NewTeamMemberID(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		Status:    TeamMemberStatusPending,
		InvitedAt: time.Now(),
	}
}

// =============================================================================
// TeamMember Methods
// =============================================================================

// FullName returns the team member's full name.
func (t *TeamMember) FullName() string {
	return t.FirstName + " " + t.LastName
}

// Accept accepts the invitation and links the user account.
func (t *TeamMember) Accept(userID UserID) {
	now := time.Now()
	t.UserID = &userID
	t.Status = TeamMemberStatusActive
	t.JoinedAt = &now
}

// Deactivate deactivates the team member.
func (t *TeamMember) Deactivate() {
	t.Status = TeamMemberStatusInactive
}

// Activate activates an inactive team member.
func (t *TeamMember) Activate() {
	t.Status = TeamMemberStatusActive
}

// UpdateRole updates the team member's role.
func (t *TeamMember) UpdateRole(role TeamRole) {
	t.Role = role
}

// UpdateInfo updates the team member's info.
func (t *TeamMember) UpdateInfo(firstName, lastName string) {
	t.FirstName = firstName
	t.LastName = lastName
}

// IsActive returns true if the team member is active.
func (t *TeamMember) IsActive() bool {
	return t.Status == TeamMemberStatusActive
}

// IsPending returns true if the invitation is pending.
func (t *TeamMember) IsPending() bool {
	return t.Status == TeamMemberStatusPending
}

// HasPermission checks if the team member has a specific permission.
func (t *TeamMember) HasPermission(permission Permission) bool {
	return t.Role.HasPermission(permission)
}

// =============================================================================
// TeamRole Enum
// =============================================================================

// TeamRole represents the role of a team member.
type TeamRole string

const (
	TeamRoleAdmin   TeamRole = "admin"   // Full access
	TeamRoleManager TeamRole = "manager" // Manage offers and bookings
	TeamRoleStaff   TeamRole = "staff"   // View and check-in only
	TeamRoleViewer  TeamRole = "viewer"  // Read-only access
)

// AllTeamRoles returns all valid team roles.
func AllTeamRoles() []TeamRole {
	return []TeamRole{
		TeamRoleAdmin,
		TeamRoleManager,
		TeamRoleStaff,
		TeamRoleViewer,
	}
}

// IsValid checks if the role is valid.
func (r TeamRole) IsValid() bool {
	for _, valid := range AllTeamRoles() {
		if r == valid {
			return true
		}
	}
	return false
}

// String returns the string representation.
func (r TeamRole) String() string {
	return string(r)
}

// =============================================================================
// Permission Enum
// =============================================================================

// Permission represents a specific permission.
type Permission string

const (
	PermissionViewDashboard        Permission = "view_dashboard"
	PermissionManageOffers         Permission = "manage_offers"
	PermissionManageEstablishments Permission = "manage_establishments"
	PermissionManageTeam           Permission = "manage_team"
	PermissionViewBookings         Permission = "view_bookings"
	PermissionCheckIn              Permission = "check_in"
	PermissionViewAnalytics        Permission = "view_analytics"
	PermissionManageSettings       Permission = "manage_settings"
)

// HasPermission checks if the role has a specific permission.
func (r TeamRole) HasPermission(perm Permission) bool {
	permissions := r.Permissions()
	for _, p := range permissions {
		if p == perm {
			return true
		}
	}
	return false
}

// Permissions returns all permissions for the role.
func (r TeamRole) Permissions() []Permission {
	switch r {
	case TeamRoleAdmin:
		return []Permission{
			PermissionViewDashboard,
			PermissionManageOffers,
			PermissionManageEstablishments,
			PermissionManageTeam,
			PermissionViewBookings,
			PermissionCheckIn,
			PermissionViewAnalytics,
			PermissionManageSettings,
		}
	case TeamRoleManager:
		return []Permission{
			PermissionViewDashboard,
			PermissionManageOffers,
			PermissionViewBookings,
			PermissionCheckIn,
			PermissionViewAnalytics,
		}
	case TeamRoleStaff:
		return []Permission{
			PermissionViewDashboard,
			PermissionViewBookings,
			PermissionCheckIn,
		}
	case TeamRoleViewer:
		return []Permission{
			PermissionViewDashboard,
			PermissionViewBookings,
			PermissionViewAnalytics,
		}
	default:
		return []Permission{}
	}
}

// =============================================================================
// TeamMemberStatus Enum
// =============================================================================

// TeamMemberStatus represents the status of a team member.
type TeamMemberStatus string

const (
	TeamMemberStatusPending  TeamMemberStatus = "pending"
	TeamMemberStatusActive   TeamMemberStatus = "active"
	TeamMemberStatusInactive TeamMemberStatus = "inactive"
)

// IsValid checks if the status is valid.
func (s TeamMemberStatus) IsValid() bool {
	return s == TeamMemberStatusPending ||
		s == TeamMemberStatusActive ||
		s == TeamMemberStatusInactive
}

// String returns the string representation.
func (s TeamMemberStatus) String() string {
	return string(s)
}
