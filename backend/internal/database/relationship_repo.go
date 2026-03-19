package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/atharvadani9/geofence-alert-app/internal/models"
	"github.com/jackc/pgx/v5"
)

type RelationshipRepo struct {
	db *DB
}

func NewRelationshipRepo(db *DB) *RelationshipRepo {
	return &RelationshipRepo{db: db}
}

func (r *RelationshipRepo) CreateInvite(ctx context.Context, caregiverID, trackedUserID string) (*models.Relationship, error) {
	query := `
		INSERT INTO relationships (caregiver_id, tracked_user_id, status)
		VALUES ($1, $2, 'pending')
		RETURNING id, caregiver_id, tracked_user_id, status, created_at, updated_at
	`
	relationship := &models.Relationship{}
	err := r.db.Pool().QueryRow(ctx, query, caregiverID, trackedUserID).Scan(
		&relationship.ID,
		&relationship.CaregiverID,
		&relationship.TrackedUserID,
		&relationship.Status,
		&relationship.CreatedAt,
		&relationship.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create invite: %w", err)
	}

	return relationship, nil
}

func (r *RelationshipRepo) GetRelationship(ctx context.Context, caregiverID, trackedUserID string) (*models.Relationship, error) {
	query := `
		SELECT id, caregiver_id, tracked_user_id, status, created_at, updated_at
		FROM relationships
		WHERE caregiver_id = $1 AND tracked_user_id = $2
	`
	relationship := &models.Relationship{}
	err := r.db.Pool().QueryRow(ctx, query, caregiverID, trackedUserID).Scan(
		&relationship.ID,
		&relationship.CaregiverID,
		&relationship.TrackedUserID,
		&relationship.Status,
		&relationship.CreatedAt,
		&relationship.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("relationship not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get relationship: %w", err)
	}

	return relationship, nil
}

func (r *RelationshipRepo) AcceptInvite(ctx context.Context, relationshipID string) error {
	query := `
		UPDATE relationships
		SET status = 'accepted', updated_at = now()
		WHERE id = $1
	`
	_, err := r.db.Pool().Exec(ctx, query, relationshipID)
	return err
}

func (r *RelationshipRepo) RejectInvite(ctx context.Context, relationshipID string) error {
	query := `
		UPDATE relationships
		SET status = 'rejected', updated_at = now()
		WHERE id = $1
	`
	_, err := r.db.Pool().Exec(ctx, query, relationshipID)
	return err
}

func (r *RelationshipRepo) DeleteRelationship(ctx context.Context, relationshipID string) error {
	query := `
		DELETE FROM relationships
		WHERE id = $1
	`
	_, err := r.db.Pool().Exec(ctx, query, relationshipID)
	return err
}

func (r *RelationshipRepo) GetCaregivers(ctx context.Context, trackedUserID string) ([]*models.User, error) {
	query := `
		SELECT u.id, u.email, u.name, u.role, u.created_at, u.updated_at
		FROM users u
		JOIN relationships r ON u.id = r.caregiver_id
		WHERE r.tracked_user_id = $1 AND r.status = 'accepted'
	`
	rows, err := r.db.Pool().Query(ctx, query, trackedUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get caregivers: %w", err)
	}
	defer rows.Close()

	caregivers := []*models.User{}
	for rows.Next() {
		caregiver := &models.User{}
		err := rows.Scan(
			&caregiver.ID,
			&caregiver.Email,
			&caregiver.Name,
			&caregiver.Role,
			&caregiver.CreatedAt,
			&caregiver.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan caregiver: %w", err)
		}
		caregivers = append(caregivers, caregiver)
	}

	return caregivers, nil
}

func (r *RelationshipRepo) GetTrackedUsers(ctx context.Context, caregiverID string) ([]*models.User, error) {
	query := `
		SELECT u.id, u.email, u.name, u.role, u.created_at, u.updated_at
		FROM users u
		JOIN relationships r ON u.id = r.tracked_user_id
		WHERE r.caregiver_id = $1 AND r.status = 'accepted'
	`
	rows, err := r.db.Pool().Query(ctx, query, caregiverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracked users: %w", err)
	}
	defer rows.Close()

	trackedUsers := []*models.User{}
	for rows.Next() {
		trackedUser := &models.User{}
		err := rows.Scan(
			&trackedUser.ID,
			&trackedUser.Email,
			&trackedUser.Name,
			&trackedUser.Role,
			&trackedUser.CreatedAt,
			&trackedUser.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tracked user: %w", err)
		}
		trackedUsers = append(trackedUsers, trackedUser)
	}

	return trackedUsers, nil
}

func (r *RelationshipRepo) CheckRelationshipExists(ctx context.Context, caregiverID, trackedUserID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM relationships
			WHERE caregiver_id = $1 AND tracked_user_id = $2
		)
	`
	var exists bool
	err := r.db.Pool().QueryRow(ctx, query, caregiverID, trackedUserID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check relationship: %w", err)
	}

	return exists, nil
}
