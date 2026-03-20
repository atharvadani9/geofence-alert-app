package database

import (
	"context"

	"github.com/atharvadani9/geofence-alert-app/internal/models"
)

type GeofenceRepo struct {
	db *DB
}

func NewGeofenceRepo(db *DB) *GeofenceRepo {
	return &GeofenceRepo{db: db}
}

func (r *GeofenceRepo) CreateGeofence(ctx context.Context, req *models.CreateGeofenceRequest) (*models.Geofence, error) {
	query := `
		INSERT INTO geofences (name, caregiver_id, tracked_user_id, radius, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, caregiver_id, tracked_user_id, radius, latitude, longitude, is_active, created_at, updated_at
	`
	geofence := &models.Geofence{}
	err := r.db.Pool().QueryRow(ctx, query, req.Name, req.CaregiverID, req.TrackedUserID, req.Radius, req.Latitude, req.Longitude).Scan(
		&geofence.ID,
		&geofence.Name,
		&geofence.CaregiverID,
		&geofence.TrackedUserID,
		&geofence.Radius,
		&geofence.Latitude,
		&geofence.Longitude,
		&geofence.IsActive,
		&geofence.CreatedAt,
		&geofence.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return geofence, nil
}

func (r *GeofenceRepo) GetGeofenceByID(ctx context.Context, geofenceID string) (*models.Geofence, error) {
	return r.getGeofence(ctx, geofenceID, "")
}

func (r *GeofenceRepo) GetGeofenceByName(ctx context.Context, name string) (*models.Geofence, error) {
	return r.getGeofence(ctx, "", name)
}

func (r *GeofenceRepo) getGeofence(ctx context.Context, geofenceID string, name string) (*models.Geofence, error) {
	query := `
		SELECT id, name, caregiver_id, tracked_user_id, radius, latitude, longitude, is_active, created_at, updated_at
		FROM geofences
		WHERE id = $1 OR name = $2
	`
	geofence := &models.Geofence{}
	err := r.db.Pool().QueryRow(ctx, query, geofenceID, name).Scan(
		&geofence.ID,
		&geofence.Name,
		&geofence.CaregiverID,
		&geofence.TrackedUserID,
		&geofence.Radius,
		&geofence.Latitude,
		&geofence.Longitude,
		&geofence.IsActive,
		&geofence.CreatedAt,
		&geofence.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return geofence, nil
}

func (r *GeofenceRepo) ListGeofencesForUser(ctx context.Context, req *models.ListGeofencesForUserRequest) ([]*models.Geofence, error) {
	query := `
		SELECT id, name, caregiver_id, tracked_user_id, radius, latitude, longitude, is_active, created_at, updated_at
		FROM geofences
		WHERE (caregiver_id = $1 AND $2 = 'caregiver') OR (tracked_user_id = $1 AND $2 = 'tracked')
		AND (is_active = true AND $3 = 'active') OR (is_active = false AND $3 = 'inactive')
		OR ($3 = 'all')
	`
	rows, err := r.db.Pool().Query(ctx, query, req.UserID, req.Role, req.Filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	geofences := []*models.Geofence{}
	for rows.Next() {
		geofence := &models.Geofence{}
		err := rows.Scan(
			&geofence.ID,
			&geofence.Name,
			&geofence.CaregiverID,
			&geofence.TrackedUserID,
			&geofence.Radius,
			&geofence.Latitude,
			&geofence.Longitude,
			&geofence.IsActive,
			&geofence.CreatedAt,
			&geofence.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		geofences = append(geofences, geofence)
	}
	return geofences, nil
}

func (r *GeofenceRepo) UpdateGeofence(ctx context.Context, geofence *models.Geofence) error {
	query := `
		UPDATE geofences
		SET name = $1, radius = $2, latitude = $3, longitude = $4, is_active = $5, updated_at = now()
		WHERE id = $6
	`
	_, err := r.db.Pool().Exec(ctx, query, geofence.Name, geofence.Radius, geofence.Latitude, geofence.Longitude, geofence.IsActive, geofence.ID)
	return err
}

func (r *GeofenceRepo) DeleteGeofence(ctx context.Context, geofenceID string) error {
	query := `
		DELETE FROM geofences
		WHERE id = $1
	`
	_, err := r.db.Pool().Exec(ctx, query, geofenceID)
	return err
}

func (r *GeofenceRepo) CheckGeofenceOwnership(ctx context.Context, geofenceID, caregiverID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM geofences
			WHERE id = $1 AND caregiver_id = $2
		)
	`
	var exists bool
	err := r.db.Pool().QueryRow(ctx, query, geofenceID, caregiverID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
