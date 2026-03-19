package models

import "time"

type Relationship struct {
	ID            string    `json:"id"`
	CaregiverID   string    `json:"caregiver_id"`
	TrackedUserID string    `json:"tracked_user_id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type InviteRequest struct {
	TrackedUserEmail string `json:"tracked_user_email" validate:"required,email"`
}

type UpdateInviteRequest struct {
	RelationshipID string `json:"relationship_id" validate:"required"`
	Status         string `json:"status" validate:"required,oneof=accepted rejected"`
}

type DeleteRelationshipRequest struct {
	RelationshipID string `json:"relationship_id" validate:"required"`
}
